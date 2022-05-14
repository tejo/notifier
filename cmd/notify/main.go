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

var serverUrl string
var interval time.Duration
var runDummyServer bool

func main() {
	setupFlags()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	if runDummyServer {
		go dummyServer()
		time.Sleep(time.Second * 1)
	}

	processInput(shutdown)
}

func processInput(shutdown chan os.Signal) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		for {
			select {
			case <-shutdown:
				log.Println("Gracefully exiting")
				// do stuff before exiting
				return
			case <-time.After(interval):
				notifier.Notify(serverUrl, scanner.Text())
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
		time.Sleep(time.Second * 5)
		body, _ := ioutil.ReadAll(r.Body)
		log.Printf("Received request: %s", body)
		w.Write([]byte("OK"))
	})
	log.Println("Starting dummy server")
	http.ListenAndServe(":8080", nil)
}
