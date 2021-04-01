package starter

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"sync/atomic"
	config2 "thingworks.net/thingworks/common/autoconfig/config"
	"thingworks.net/thingworks/common/https"
	"thingworks.net/thingworks/common/starter/service"
	"thingworks.net/thingworks/common/utils/strings2"
)

var wg sync.WaitGroup
var routineCounter uint32 = 0

type defaultStarter struct {
	gateway  https.GatewayAdaptor
	services []ServiceStarter
}

func (starter *defaultStarter) RegisterStarter(service ServiceStarter) {
	starter.services = append(starter.services, service)
}

func (starter *defaultStarter) Run([]string) {
	port := strings2.Concat(":", strings2.Itoa(starter.port()))

	if config2.DefaultConfig().Log.Debug {
		log.SetLevel(log.DebugLevel)
	}

	starter.StartAllServices()

	wg.Add(1)
	plusRoutineCounter()

	go starter.ListenAndServe(port)

	wg.Wait()
}

func (starter *defaultStarter) RegisterResource(resourceMap https.ResourceMap) ApplicationStarter {
	starter.gateway.RegisterResource(resourceMap)
	return starter
}

func (starter *defaultStarter) Stop() {
	starter.gateway.Close()

	for count := 0; count < int(routineCounter); count++ {
		wg.Done()
	}

	resetRoutineCounter()
}

func resetRoutineCounter() bool {
	return atomic.CompareAndSwapUint32(&routineCounter, routineCounter, 0)
}

func (starter *defaultStarter) StartAllServices() {
	for _, serviceStarter := range starter.services {
		err := serviceStarter.Start()
		if err != nil {
			log.WithFields(log.Fields{
				"config": config2.DefaultConfig(),
			}).Errorf("Exception when start service %s", serviceStarter.ServiceName())
		}
	}
}

func (starter *defaultStarter) port() int {
	return config2.DefaultConfig().App.Port
}

func (starter *defaultStarter) ListenAndServe(port string) {
	http.Handle("/", starter.gateway)
	log.Infof("Service start at port %s", port)

	wg.Add(1)
	plusRoutineCounter()

	go starter.gateway.Start()

	err := http.ListenAndServe(port, nil)

	if err != nil {
		panic(ApplicationStartError{err: err})
	}
}

func plusRoutineCounter() bool {
	return atomic.CompareAndSwapUint32(&routineCounter, routineCounter, atomic.AddUint32(&routineCounter, 1))
}

func GetDefaultAppStarter(opts ConfigOptions) ApplicationStarter {
	config2.Init(config2.AppArgs{
		ConfigLocation: &opts.Path,
	})

	starter := &defaultStarter{
		gateway: https.NewGateway(),
	}

	starter.RegisterStarter(service.NewMongoStarter())
	starter.RegisterStarter(service.NewMqttServiceStarter())

	return starter
}
