package notifier

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/alitto/pond"
	log "github.com/sirupsen/logrus"
)

var (
	pool       *pond.WorkerPool
	maxWorkers = 100
	maxQueue   = 1000
)

type API struct {
	Client *http.Client
	URL    string
}

func New(URL string) *API {
	pool = pond.New(maxWorkers, maxQueue)
	return &API{&http.Client{}, URL}
}

func (api *API) Notify(message string) {
	pool.Submit(func() {
		_, err := api.PostMessage(message)
		if err != nil {
			log.Errorf("Error posting message: %s", err)
		}
	})
}

func (api *API) StopWorkers() {
	log.Println("Wait enqueued jobs to finish")
	pool.StopAndWait()
	log.Println("All jobs finished, exiting.")
	os.Exit(0)
}

func (api *API) PostMessage(message string) ([]byte, error) {
	resp, err := api.Client.Post(api.URL, "application/x-www-form-urlencoded", strings.NewReader(message))
	if err != nil {
		log.Errorf("Error posting message to %s: %s", api.URL, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Error posting message to %s: %s", api.URL, resp.Status)
		return nil, fmt.Errorf("error posting message to %s: %s", api.URL, resp.Status)
	}

	log.Printf("Posted message to %s", api.URL)

	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
