package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"notifier/pkg/notifier"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	serverUrl      string
	interval       time.Duration
	shutdown       chan os.Signal
	runDummyServer bool
)

func main() {
	setupFlags()

	shutdown = make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	if runDummyServer {
		go dummyServer()
	}

	processInput()
}

func processInput() {
	scanner := bufio.NewScanner(os.Stdin)
	notifier := notifier.New(serverUrl)

	for {
		select {
		case <-shutdown:
			notifier.StopWorkers()
			return
		case <-time.After(interval):
			if scanner.Scan() {
				notifier.Notify(scanner.Text())
			} else {
				notifier.StopWorkers()
			}
		}
	}
}

func setupFlags() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s --url=URL [<flags>]\n\nflags:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&serverUrl, "url", "http://localhost:8080/notify", "URL of the server to post to")
	flag.DurationVar(&interval, "interval", time.Second*5, "interval to wait between posts")
	flag.BoolVar(&runDummyServer, "dummyserver", false, "runs a dummy server for testing")
	flag.Parse()
}

func dummyServer() {
	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 3)
		body, _ := ioutil.ReadAll(r.Body)
		log.Printf("Received request: %s", body)
		w.Write([]byte("OK"))
	})
	log.Println("Starting dummy server")
	http.ListenAndServe(":8080", nil)
}
