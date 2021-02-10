package config

import (
	"github.com/thingworks/common/utils/strings2"
	"time"
)

const mongoUriPrefix = "mongodb://"

type MongoConfig struct {
	Uri      string
	Host     string
	Port     string
	Username string
	Password string
	DataBase string
	Timeout  time.Duration
}

func (config MongoConfig) GetUri() string {
	if strings2.IsNotBlank(config.Uri) {
		return config.Uri
	}

	return strings2.Concat(mongoUriPrefix, strings2.Join([]string{config.Host, config.Port}, ":"))
}
