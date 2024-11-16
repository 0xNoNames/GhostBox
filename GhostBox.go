package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/0xNoNames/GhostBox/client"
	"github.com/0xNoNames/GhostBox/utils"
	"github.com/radovskyb/watcher"
)

var (
	input       *string
	output      *string
	helpFlag    *bool
	downloadDir string
	watchedDir  string
)

func init() {
	input = flag.String("i", "", "Watched directory (where the .torrent files will be added)")
	output = flag.String("o", "", "Output directory (where the downloaded files will be stored)")
	helpFlag = flag.Bool("help", false, "Show this message")

	// Define how to use the program
	flag.Usage = func() {
		fmt.Printf("GhostBox - Watcher torrent client\n\n")
		fmt.Printf("Usage: %s OPTIONS\n\n", os.Args[0])
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}

	// Parse command-line flags
	flag.Parse()

	// If -help flag is set, show the usage and exit
	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	// Use the value of -download-dir flag
	downloadDir = *output
	if downloadDir == "" {
		log.Print("Output download directory must be set.")
		flag.Usage()
		os.Exit(1)
	} else if !utils.Exist(downloadDir) {
		log.Fatal("Output download directory must exist.")
	}

	watchedDir = *input
	if watchedDir == "" {
		log.Print("Watched directory must be set.")
		flag.Usage()
		os.Exit(1)
	} else if !utils.Exist(watchedDir) {
		log.Fatal("Watched directory must exist.")
	}
}

func main() {
	torrentClient, err := client.New(downloadDir)
	if err != nil {
		panic("Error creating torrent client")
	}
	err = torrentClient.Start()
	if err != nil {
		panic("Error starting torrent client")
	}

	w := watcher.New()

	// Only notify rename and move events.
	w.FilterOps(watcher.Move, watcher.Create, watcher.Remove)

	// Only files that match the regular expression during file listings will be watched.
	r := regexp.MustCompile(`.+\.torrent$`)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				// If a torrent file is added to the watched directory, adds it to the torrent client
				if event.Op == watcher.Create && r.MatchString(event.Name()) {
					log.Printf("New file matching detected: \"%s\"", event.Name())
					torrentClient.AddTorrent(event.Path)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this dir for changes.
	if err := w.Add(watchedDir); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and dirs currently being watched and their paths.
	for path := range w.WatchedFiles() {
		log.Printf("Watching directory for changes: \"%s\"", path)
		// If a torrent file is already in the watched directory, adds it to the torrent client
		if r.MatchString(path) {
			log.Printf("New file matching detected: \"%s\"", path)
			torrentClient.AddTorrent(path)
		}
	}

	// Start the watching process - it'll check for changes every 5 seconds.
	if err := w.Start(time.Second * 5); err != nil {
		log.Fatalln(err)
	}
}
