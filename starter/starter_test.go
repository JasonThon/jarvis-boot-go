package starter

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"thingworks/common/https"
	"thingworks/common/utils/strings2"
	"time"
)

const greeting = "hello"

type testResource struct{}

func (t testResource) Handlers() https.HandlerMap {
	return https.HandlerMap{
		"hello": https.Get(t.hello),
		"": https.Get(t.overview),
	}
}

func request(path string) *http.Response {
	resp, err := http.Get("http://localhost:9090/" + path)

	if err != nil {
		log.Fatalf("Error happens: %s", err.Error())
	}

	return resp
}

func (t testResource) hello(w http.ResponseWriter, _ *https.HttpRequest) {
	_, err := w.Write(strings2.ToByte(greeting))

	if err != nil {
		log.Fatalf("Error happens when writing data: %v", err)
	}
}

func (t testResource) overview(writer http.ResponseWriter, _ *https.HttpRequest) {
	writer.Write(strings2.ToByte("ok"))
}

func TestDefaultStarter_Run(t *testing.T) {
	s := GetDefaultAppStarter(ConfigOptions{Path: "./config.yaml"})
	s.RegisterResource(https.ResourceMap{
		"/teststarter": testResource{},
	})
	s.Run(nil)
	time.Sleep(3 * time.Second)

	resp := request("/teststarter/hello")
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Error happens when start server: %v", r)
		}
	}()

	assert.True(t, resp.StatusCode == 200)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error happens: %v", err)
	}

	respData := string(body)

	assert.Equal(t, greeting, respData)

	resp = request("/teststarter")

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error happens: %v", err)
	}

	respData = string(body)

	assert.Equal(t, "ok", respData)
}
