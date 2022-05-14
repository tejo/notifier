package main

import (
	"flag"
	"fmt"
	"notifier/pkg/notifier"
	"time"
)

func main() {
	var serverUrl string
	var interval time.Duration
	var runDummyServer bool

	flag.StringVar(&serverUrl, "url", "", "URL of the server to post to")
	flag.DurationVar(&interval, "interval", time.Second*5, "interval to wait between posts")
	flag.BoolVar(&runDummyServer, "dummyserver", false, "runs a dummy server for testing")
	flag.Parse()

	fmt.Print(runDummyServer)
	notifier.Notify("http://127.0.0.1:4567", "Hello, world!")
}
