package config

import (
	"thingworks.net/thingworks/jarvis-boot/utils/strings2"
	"time"
)

const mongoUriPrefix = "mongodb://"

type MongoConfig struct {
	Uri           string
	Host          string
	Port          string
	Username      string
	Password      string
	DataBase      string
	Timeout       time.Duration
	AuthMechanism string
}

func (config MongoConfig) GetUri() string {
	if strings2.IsNotBlank(config.Uri) {
		return config.Uri
	}

	return strings2.Concat(mongoUriPrefix, strings2.Join([]string{config.Host, config.Port}, ":"))
}

func (config MongoConfig) IsValidConf() bool {
	return config.isValidHost() && config.isValidDb()
}

func (config MongoConfig) isValidHost() bool {
	return strings2.IsNotBlank(config.Uri) || (strings2.IsNotBlank(config.Port) && strings2.IsNotBlank(config.Host))
}

func (config MongoConfig) isValidDb() bool {
	return strings2.IsNotBlank(config.DataBase)
}
