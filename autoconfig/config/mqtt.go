package config

import (
	"thingworks.net/thingworks/jarvis-boot/utils/strings2"
	"time"
)

const defaultRetry = 3

type MqttConfig struct {
	Host              string
	Port              int
	Broker            string
	ClientId          string        `yaml:"clientId"`
	KeepAlive         time.Duration `yaml:"keepAlive"`
	PingTimeout       time.Duration `yaml:"pingTimeout"`
	ConnectionTimeout time.Duration `yaml:"connectionTimeout"`
	Username          string
	Password          string
	Retry             int
}

func (mqtt MqttConfig) GetBroker() string {
	if strings2.IsNotBlank(mqtt.Broker) {
		return mqtt.Broker
	}

	host := strings2.Split(mqtt.Host, ":")[0]

	return strings2.Concat(host, ":", strings2.Itoa(mqtt.Port))
}

func (mqtt MqttConfig) IsValid() bool {
	return strings2.IsNotBlank(mqtt.GetBroker()) &&
		((mqtt.Port > 0 && strings2.IsNotBlank(mqtt.Host)) ||
			strings2.IsNotBlank(mqtt.Broker))
}

func (mqtt MqttConfig) GetRetry() int {
	if mqtt.Retry <= 0 {
		return defaultRetry
	}

	return mqtt.Retry
}
