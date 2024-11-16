package client

import "github.com/anacrolix/log"

type Information struct {
	Name         string
	Progress     string
	Seeders      string
	Leechers     string
	Torrentspeed string
	ETA          string
}

type TorrentInformation struct {
	Information Information
	index       int
	finished    bool
	aborted     bool
	dropped     bool
	path        string
}

func defaultTorrentInformation(name string, index int, path string) TorrentInformation {
	return TorrentInformation{
		Information: Information{
			Name:         name,
			Progress:     "",
			Seeders:      "",
			Leechers:     "",
			Torrentspeed: "",
			ETA:          "",
		},
		index:    index,
		finished: false,
		aborted:  false,
		dropped:  false,
		path:     path,
	}
}

func (d *TorrentInformation) Abort() {
	d.aborted = true
	if !d.dropped {
		log.Printf("Aborted torrent: \"%s\"", d.Information.Name)
		d.Dropped()
	}
}

func (d *TorrentInformation) Dropped() {
	if !d.dropped {
		d.dropped = true
		d.Information.Seeders = "0"
		d.Information.Leechers = "0"
		d.Information.Torrentspeed = "0.00MB/s"
	}
}

func (i *Information) logTorrentInfo() {
	log.Printf("[TRACK] %s | PROGRESS: %s | SEED: %s | LEECH: %s | SPEED: %s | ETA: %s",
		i.Name, i.Progress, i.Seeders, i.Leechers, i.Torrentspeed, i.ETA)
}
