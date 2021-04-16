package config

import (
	"thingworks.net/thingworks/jarvis-boot/utils/strings2"
	"time"
)

type MqttConfig struct {
	Host              string
	Port              int
	ClientId          string        `yaml:"clientId"`
	KeepAlive         time.Duration `yaml:"keepAlive"`
	PingTimeout       time.Duration `yaml:"pingTimeout"`
	ConnectionTimeout time.Duration `yaml:"connectionTimeout"`
	Username          string
	Password          string
}

func (mqtt MqttConfig) Broker() string {
	return strings2.Concat(mqtt.Host, ":", strings2.Itoa(mqtt.Port))
}

func (mqtt MqttConfig) IsValid() bool {
	return strings2.IsNotBlank(mqtt.Broker()) && mqtt.Port > 0 && strings2.IsNotBlank(mqtt.Host)
}
