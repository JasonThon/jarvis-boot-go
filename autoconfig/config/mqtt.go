package config

import (
	"thingworks.net/thingworks/jarvis-boot/utils/strings2"
	"time"
)

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
