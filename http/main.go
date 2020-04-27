package main

import (
	"log"

	"github.com/hanwen/go-fuse/v2/fs"
)

func main() {
	urls := []string{
		"http://localhost:8080/chinook.db",
		"http://localhost:8080/enwiki.db",
		"https://dbhub.io/x/download/justinclift/backblaze-drive_stats.db?commit=5e9fe0d80b7e85bb3b027499c336b4f4b8c071a5a55075aab06adbab50c083c9",
	}

	opts := &fs.Options{}
	opts.Debug = false

	r := &HttpRoot{
		URLs: urls,
	}
	server, err := fs.Mount("/tmp/distribyted", r, opts)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}

	defer server.Unmount()
	server.Wait()
}
