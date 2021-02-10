package config

import (
	"errors"
	"github.com/thingworks/common/utils/strings2"
	"sync"
)

const (
	defaultConfigFileName = "config.yaml"
)

type AppArgs struct {
	ConfigLocation   *string
	PoliciesLocation *string
}

var config *AppConfig
var once sync.Once

type RefresherConfig struct {
	Path string
}

type LogConfig struct {
	Debug bool
}

type ServerConfig struct {
	Port int
	Name string
}

func (s ServerConfig) Check() error {
	if s.Port <= 0 {
		return errors.New("Invalid Prt: " + strings2.Itoa(s.Port))
	}

	return nil
}

type RegistryConfig struct {
	Recover RecoverConfig
}
