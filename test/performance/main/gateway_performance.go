package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"os"
	"thingworks.net/thingworks/jarvis-boot/https"
	"thingworks.net/thingworks/jarvis-boot/starter"
	"thingworks.net/thingworks/jarvis-boot/utils/strings2"
)

const greet = "hello"

type purseudoResource struct{}

func (t *purseudoResource) Handlers() https.HandlerMap {
	return https.HandlerMap{
		"hello": https.Get(t.hello),
	}
}

func (t *purseudoResource) hello(writer http.ResponseWriter, request *https.HttpRequest) {
	bytes := strings2.ToByte(greet)

	_, err := writer.Write(bytes)

	if err != nil {
		log.Fatalf("Error happens when writing data: %v", err)
	}
}

func main() {
	dir, _ := os.Getwd()
	log.Info(dir)

	s := starter.GetDefaultAppStarter(starter.ConfigOptions{Path: strings2.Concat(dir, "/performance/config.yaml")})
	s.RegisterResource(https.ResourceMap{
		"/purseudo": &purseudoResource{},
	})

	// 性能监控用，访问 http://localhost:8080/debug/pprof/ 来查看资源消耗
	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()

	s.Run(nil)

}
