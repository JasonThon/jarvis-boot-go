package service

import (
	"thingworks/common/autoconfig/config"
	"thingworks/common/autoconfig/mongodb"
)

type MongoServiceStarter struct{}

func (m *MongoServiceStarter) Start(configs config.AppConfig) error {
	return mongodb.InitMongoTemplate(configs.Mongodb)
}

func (m *MongoServiceStarter) ServiceName() string {
	return "mongodb"
}

func NewMongoStarter() *MongoServiceStarter {
	return &MongoServiceStarter{}
}

