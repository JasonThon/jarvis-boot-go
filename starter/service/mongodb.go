package service

import (
	"thingworks.net/thingworks/common/autoconfig/config"
	"thingworks.net/thingworks/common/autoconfig/mongodb"
)

type MongoServiceStarter struct{}

func (m *MongoServiceStarter) Start() error {
	return mongodb.InitMongoTemplate(config.DefaultConfig().Mongodb)
}

func (m *MongoServiceStarter) ServiceName() string {
	return "mongodb"
}

func NewMongoStarter() *MongoServiceStarter {
	return &MongoServiceStarter{}
}
