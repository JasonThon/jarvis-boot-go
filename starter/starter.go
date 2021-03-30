package starter

import (
	"fmt"
	"thingworks.net/thingworks/common/autoconfig/config"
	"thingworks.net/thingworks/common/https"
)

type ApplicationStarter interface {
	Run(args []string)
	StartAllServices()
	RegisterStarter(service ServiceStarter)
	RegisterResource(https.ResourceMap) ApplicationStarter
	Stop()
}

type ServiceStarter interface {
	Start(appConfig config.AppConfig) error
	ServiceName() string
}

type ConfigOptions struct {
	Path string
}

type ApplicationStartError struct {
	err error
}

func (appErr *ApplicationStartError) Error() string {
	return fmt.Sprintf("Exception happens when application starts: %v", appErr.err)
}
