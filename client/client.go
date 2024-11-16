package client

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/0xNoNames/GhostBox/utils"
	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
)

// Constants
const (
	megabyte           = 1024 * 1024
	torrentSpeedFormat = "%.2fMB/s"
)

// Model represents the torrent client model.
type Model struct {
	Torrents []TorrentInformation // Slice to store information about each torrent being downloaded.
	client   *torrent.Client      // The underlying torrent client used to manage torrents.
	files    []string             // Paths to the torrent files to be downloaded.
	mux      sync.Mutex           // Mutex to protect concurrent access to Torrents slice.
}

// New creates a new torrent client model.
func New(downloadDir string) (*Model, error) {
	// Create a new configuration for the torrent client.
	clientConfig := createClientConfig(downloadDir)
	// Create a new torrent client.
	client, err := torrent.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}
	TorrentsInformation := make([]TorrentInformation, 0)
	// Create a new model.
	return &Model{Torrents: TorrentsInformation, client: client}, nil
}

// AddTorrent adds a new torrent info to the model.
func (m *Model) AddTorrent(path string) error {
	var torrentInformation TorrentInformation

	// Check if the torrent is already added
	if m.isTorrentAdded(path) {
		return fmt.Errorf("Torrent already added: %s", path)
	}

	m.mux.Lock()         // Lock the mutex before modifying the Torrents slice.
	defer m.mux.Unlock() // Ensure we unlock the mutex even if there's a panic.

	length := len(m.Torrents)

	t, err := m.client.AddTorrentFromFile(path)
	if err != nil {
		torrentInformation = defaultTorrentInformation(err.Error(), length, path)
		torrentInformation.Abort()
	} else {
		torrentInformation = defaultTorrentInformation(t.Info().Name, length, path)
	}
	// Add the torrent info to the model.
	m.Torrents = append(m.Torrents, torrentInformation)

	// Start downloading the torrent.
	t.DownloadAll()

	// Start tracking the torrent in a separate goroutine.
	go m.trackTorrent(t, length)

	return nil
}

// isTorrentAdded checks if a torrent with the given path is already added.
func (m *Model) isTorrentAdded(path string) bool {
	for _, torrentInformation := range m.Torrents {
		if torrentInformation.path == path {
			return true
		}
	}
	return false
}

// Start starts the download process for the client.
func (m *Model) Start() error {
	for _, path := range m.files {
		err := m.AddTorrent(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// trackTorrent tracks the download progress of a torrent.
func (m *Model) trackTorrent(t *torrent.Torrent, index int) {
	<-t.GotInfo()

	name := t.Info().Name
	startTime := time.Now()
	startSize := t.BytesCompleted()

	for {
		m.mux.Lock() // Lock the mutex before accessing the Torrents slice.
		torrentInformation := &m.Torrents[index]
		bytesCompleted := t.BytesCompleted()
		totalLength := t.Info().TotalLength()
		elapsedTime := time.Since(startTime)

		if torrentInformation.aborted {
			if !torrentInformation.dropped {
				t.Drop()
				torrentInformation.dropped = true
			}
		} else if bytesCompleted >= totalLength {
			torrentInformation.finished = true
			log.Printf("Finished torrent: \"%s\" at \"%s\"", torrentInformation.Information.Name, torrentInformation.path)
			// Delete the torrent file after finishing
			err := os.Remove(torrentInformation.path)
			if err != nil {
				log.Printf("Error deleting torrent file: %v", err)
			} else {
				log.Printf("Successfully deleted torrent file: \"%s\"", torrentInformation.path)
			}
		} else if torrentInformation.dropped {
			t, _ = m.client.AddTorrentFromFile(torrentInformation.path)
			t.DownloadAll()
			startSize = t.BytesCompleted()
			startTime = time.Now()
			torrentInformation.dropped = false
		} else if bytesCompleted < totalLength {
			remainingBytes := totalLength - bytesCompleted
			downloadRate := calculateDownloadRate(bytesCompleted, startSize, elapsedTime)
			torrentInformation.Information = Information{
				Name:         name,
				Progress:     utils.FormatBytesProgress(bytesCompleted, totalLength),
				Seeders:      strconv.Itoa(t.Stats().ConnectedSeeders),
				Leechers:     strconv.Itoa(t.Stats().ActivePeers - t.Stats().ConnectedSeeders),
				Torrentspeed: fmt.Sprintf(torrentSpeedFormat, downloadRate),
				ETA:          utils.FormatDuration(calculateETA(remainingBytes, downloadRate)),
			}
			torrentInformation.Information.logTorrentInfo()
		}

		m.mux.Unlock() // Unlock the mutex after accessing the Torrents slice.

		if torrentInformation.finished || torrentInformation.aborted {
			break
		}
		time.Sleep(5 * time.Second)
	}
}

// calculateDownloadRate calculates the download rate in MB/s.
func calculateDownloadRate(bytesCompleted, startSize int64, elapsedTime time.Duration) float64 {
	return float64(bytesCompleted-startSize) / elapsedTime.Seconds() / megabyte
}

// calculateETA calculates the estimated time of arrival for the download completion.
func calculateETA(remainingBytes int64, downloadRate float64) time.Duration {
	if downloadRate > 0 {
		return time.Duration(float64(remainingBytes)/downloadRate) / megabyte
	}
	return time.Duration(0)
}

// Abort aborts all Torrents in the client model.
func (m *Model) Abort() {
	for _, torrentInformation := range m.Torrents {
		torrentInformation.Abort()
	}
}

// getAvailablePort returns an available port by listening on a random port and extracting the chosen port.
func getAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	_, portString, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return 0, err
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		return 0, err
	}

	return port, nil
}

// createClientConfig creates a new configuration for the torrent client.
func createClientConfig(downloadDir string) *torrent.ClientConfig {
	port, err := getAvailablePort()
	if err != nil {
		return nil
	}

	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.ListenPort = port
	clientConfig.DataDir = downloadDir
	clientConfig.DisableTrackers = false
	clientConfig.Seed = false
	clientConfig.NoUpload = true
	clientConfig.DisableIPv6 = true
	clientConfig.Debug = false
	clientConfig.DisableWebtorrent = true
	clientConfig.DisableWebseeds = true
	clientConfig.DisableAcceptRateLimiting = true
	clientConfig.NoDefaultPortForwarding = true
	clientConfig.Logger.SetHandlers(log.DiscardHandler)
	return clientConfig
}
