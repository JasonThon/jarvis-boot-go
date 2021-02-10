package starter

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	config2 "github.com/thingworks/common/autoconfig/config"
	"github.com/thingworks/common/https"
	"github.com/thingworks/common/starter/service"
	"github.com/thingworks/common/utils/strings2"
)

type defaultStarter struct {
	appConfig config2.AppConfig
	gateway   *https.Gateway
	services  []ServiceStarter
}

func (starter *defaultStarter) RegisterStarter(service ServiceStarter) {
	starter.services = append(starter.services, service)
}

func (starter *defaultStarter) Run([]string) {
	port := strings2.Concat(":", strings2.Itoa(starter.port()))
	starter.StartAllServices()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go starter.ListenAndServe(port)
	wg.Wait()
}

func (starter *defaultStarter) RegisterResource(resourceMap https.ResourceMap) ApplicationStarter {
	starter.gateway.RegisterResource(resourceMap)
	return starter
}

func (starter *defaultStarter) Stop() {

}

func (starter *defaultStarter) StartAllServices() {
	for _, serviceStarter := range starter.services {
		err := serviceStarter.Start(starter.appConfig)
		if err != nil {
			log.WithFields(log.Fields{
				"config": starter.appConfig,
			}).Errorf("Exception when start service %s", serviceStarter.ServiceName())
		}
	}
}

func (starter *defaultStarter) port() int {
	return starter.appConfig.App.Port
}

func (starter *defaultStarter) ListenAndServe(port string) {
	http.Handle("/", starter.gateway)
	log.Infof("Service start at port %s", port)
	err := http.ListenAndServe(port, nil)

	if err != nil {
		panic(ApplicationStartError{err: err})
	}
}

func GetDefaultAppStarter(opts ConfigOptions) ApplicationStarter {
	config2.Init(config2.AppArgs{
		ConfigLocation: &opts.Path,
	})

	starter := &defaultStarter{
		appConfig: config2.DefaultConfig(),
		gateway:   https.NewGateway(),
	}

	starter.RegisterStarter(service.NewMongoStarter())

	return starter
}
