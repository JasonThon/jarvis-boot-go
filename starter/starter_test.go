package starter

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"thingworks.net/thingworks/common/https"
	"thingworks.net/thingworks/common/utils/strings2"
	"time"
)

const greeting = "hello"

type testResource struct{}

func (t testResource) Handlers() https.HandlerMap {
	return https.HandlerMap{
		"hello":                     https.Get(t.hello),
		"overview":                  https.Get(t.overview),
		"{id}/{containerId}/{name}": https.Post(t.returnId),
		"path/{path:.*}":            https.Get(t.path),
	}
}

func getReq(path string) *http.Response {
	resp, err := http.Get("http://localhost:9090/" + path)

	if err != nil {
		log.Fatalf("Error happens: %s", err.Error())
	}

	log.Debugf("Response is: %v", resp)

	return resp
}

func (t testResource) hello(w http.ResponseWriter, _ *https.HttpRequest) {
	bytes := strings2.ToByte(greeting)

	_, err := w.Write(bytes)

	if err != nil {
		log.Fatalf("Error happens when writing data: %v", err)
	}
}

func (t testResource) overview(writer http.ResponseWriter, _ *https.HttpRequest) {
	_, err := writer.Write(strings2.ToByte("ok"))

	if err != nil {
		log.Fatal(err)
	}
}

func (t testResource) returnId(writer http.ResponseWriter, req *https.HttpRequest) {
	id := req.GetPathParam("id")
	containerId := req.GetPathParam("containerId")
	name := req.GetPathParam("name")

	buffer := make([]byte, req.ContentLength)
	_, err2 := req.Body.Read(buffer)

	if err2 != nil {
		log.Fatal(err2)
	}

	_, err := writer.Write(strings2.ToByte(strings2.Join([]string{id, containerId, name, string(buffer)}, "-")))

	if err != nil {
		log.Fatal(err)
	}
}

func (t testResource) path(writer http.ResponseWriter, request *https.HttpRequest) {
	path := request.GetPathParam("path")

	writer.Write(strings2.ToByte(path))
}

func TestDefaultStarter_Run(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Error happens when start server: %v", r)
		}
	}()

	s := GetDefaultAppStarter(ConfigOptions{Path: "./config.yaml"})
	s.RegisterResource(https.ResourceMap{
		"/teststarter": testResource{},
	})
	go s.Run(nil)
	time.Sleep(3 * time.Second)

	var resp *http.Response

	waitgroup := sync.WaitGroup{}

	waitgroup.Add(1)
	go func() {
		respData := getResponse(t, resp, "teststarter/hello")
		assert.Equal(t, greeting, respData)
		waitgroup.Done()
	}()

	waitgroup.Add(1)
	go func() {
		respData := getResponse(t, resp, "teststarter/overview")
		assert.Equal(t, "ok", respData)
		waitgroup.Done()
	}()

	waitgroup.Add(1)
	go func() {
		respData := postResponse(t, resp, "teststarter/111/222/333")

		assert.Equal(t, "111-222-333-444", respData)
		waitgroup.Done()
	}()

	waitgroup.Add(1)
	go func() {
		respData := getResponse(t, resp, "teststarter/path/111/222/333/444")
		assert.Equal(t, "111/222/333/444", respData)
		waitgroup.Done()
	}()

	go s.Stop()
	waitgroup.Wait()
}

func postResponse(t *testing.T, resp *http.Response, path string) string {
	resp = postReq(path)

	assert.True(t, resp.StatusCode == 200)

	return extractRespBody(resp)
}

func postReq(path string) *http.Response {
	resp, err := http.Post("http://localhost:9090/"+path, "application/json", bytes.NewBuffer(strings2.ToByte("444")))

	if err != nil {
		log.Fatalf("Error happens %v", err)

		return nil
	}

	log.Debugf("Response is: %v", resp)

	return resp
}

func getResponse(t *testing.T, resp *http.Response, path string) string {
	resp = getReq(path)

	assert.True(t, resp.StatusCode == 200)

	return extractRespBody(resp)
}

func extractRespBody(resp *http.Response) string {
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error happens: %v", err)
	}

	respData := string(body)
	return respData
}
