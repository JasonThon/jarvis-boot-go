package config

import (
	"thingworks.net/thingworks/common/utils/strings2"
	"time"
)

type MqttConfig struct {
	Broker            string
	ClientId          string        `yaml:"clientId"`
	KeepAlive         time.Duration `yaml:"keepAlive"`
	PingTimeout       time.Duration `yaml:"pingTimeout"`
	ConnectionTimeout time.Duration `yaml:"connectionTimeout"`
	Username          string
	Password          string
}

func (mqtt MqttConfig) IsValid() bool {
	return strings2.IsNotBlank(mqtt.Broker)
}
