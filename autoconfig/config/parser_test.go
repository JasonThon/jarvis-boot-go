package config

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"testing"
)

var reg = regexp.MustCompile("\\$\\{(?P<env>[^}]*)\\}")

func TestRegexParser(t *testing.T) {
	assert.True(t, reg.MatchString("${env.mongodb.password}"))
	logrus.Info(reg.SubexpNames())
	logrus.Info(reg.FindStringSubmatch("${env.mongodb.password}"))
}

func TestMongoConfigParser_ParseConfig(t *testing.T) {
	mongoConf := MongoConfig{
		Uri:      "",
		Host:     "${env.mongodb.host}",
		Port:     "${env.mongodb.port}",
		Username: "${env.mongodb.username}",
		Password: "${env.mongodb.password}",
		DataBase: "${env.mongodb.schema}",
		Timeout:  10,
	}

	err := os.Setenv("env.mongodb.host", "localhost")
	logErr(err)
	err = os.Setenv("env.mongodb.port", "27017")
	logErr(err)
	logErr(os.Setenv("env.mongodb.username", ""))
	logErr(os.Setenv("env.mongodb.password", ""))
	logErr(os.Setenv("env.mongodb.schema", "jarvis"))

	parser := NewDefaultConfigParser()
	parser.ParseConfig(&mongoConf)
	assert.Equal(t, "27017", mongoConf.Port)
	assert.Equal(t, "localhost", mongoConf.Host)
	assert.Equal(t, "jarvis", mongoConf.DataBase)
}

func logErr(err error) {
	if err != nil {
		logrus.Error(err)
	}
}

func TestAppConfigParser_ParseConfig(t *testing.T) {
	mongoConf := MongoConfig{
		Uri:      "",
		Host:     "${env.mongodb.host}",
		Port:     "${env.mongodb.port}",
		Username: "${env.mongodb.username}",
		Password: "${env.mongodb.password}",
		DataBase: "${env.mongodb.schema}",
		Timeout:  10,
	}

	appConfig := AppConfig{Mongodb: mongoConf}

	err := os.Setenv("env.mongodb.host", "localhost")
	logErr(err)
	err = os.Setenv("env.mongodb.port", "27017")
	logErr(err)
	logErr(os.Setenv("env.mongodb.username", ""))
	logErr(os.Setenv("env.mongodb.password", ""))
	logErr(os.Setenv("env.mongodb.schema", "jarvis"))

	parser := AppConfigParser{NewDefaultConfigParser()}
	parser.ParseConfig(&appConfig)
	assert.Equal(t, "27017", appConfig.Mongodb.Port)
	assert.Equal(t, "localhost", appConfig.Mongodb.Host)
	assert.Equal(t, "jarvis", appConfig.Mongodb.DataBase)
}
