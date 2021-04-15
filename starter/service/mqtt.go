package service

import (
	"thingworks.net/thingworks/jarvis-boot/autoconfig/config"
	mqtt2 "thingworks.net/thingworks/jarvis-boot/autoconfig/mqtt"
)

var client *mqtt2.Client

type MqttServiceStarter struct{}

func NewMqttServiceStarter() *MqttServiceStarter {
	return &MqttServiceStarter{}
}

func (mqtt *MqttServiceStarter) Start() error {
	client = mqtt2.NewClient(config.DefaultConfig())
	mqtt2.Init(client)

	return client.Conn()
}

func (mqtt *MqttServiceStarter) ServiceName() string {
	return "mqtt"
}
