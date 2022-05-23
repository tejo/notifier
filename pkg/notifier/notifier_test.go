package notifier

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestPostMessage(t *testing.T) {
	message := "Hello, world!"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
		assert.Equal(t, req.Method, "POST")
		assert.Equal(t, body, []byte(message))
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	api := API{server.Client(), server.URL}
	body, err := api.PostMessage(message)

	assert.Nil(t, err)
	assert.Equal(t, []byte("OK"), body)
}
