package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var appParser = NewAppConfigParser()

type AppConfig struct {
	Redis     RedisConfig
	App       ServerConfig
	Log       LogConfig
	Refresher RefresherConfig
	Registry  RegistryConfig
	ApiKey    string
	Mongodb   MongoConfig
}

func (config *AppConfig) Check() error {
	return config.App.Check()
}

func DefaultConfig() AppConfig {
	if config == nil {
		panic("AppConfig file doesn NOT loaded")
	}
	return *config
}

func Init(appConfig AppArgs) AppConfig {
	once.Do(func() {
		configData, err := readConfigFile(appConfig)

		if err != nil {
			log.Fatal("Failed to load config file :" + err.Error())
		}

		parseConfigData(configData)
	})

	return *config
}

func parseConfigData(configData []byte) {
	var localConfig AppConfig
	err := yaml.Unmarshal(configData, &localConfig)

	if err != nil {
		log.Fatal("Failed to parse config file :" + err.Error())
	}

	if err = localConfig.Check(); err != nil {
		log.Fatalf("Config file is NOT valid: %s", err.Error())
	}

	appParser.ParseConfig(localConfig)

	config = &localConfig
}

func readConfigFile(appConfig AppArgs) ([]byte, error) {
	log.Println("Using config file :", *appConfig.ConfigLocation)

	var configData []byte
	var err error

	if appConfig.ConfigLocation != nil {
		configData, err = ioutil.ReadFile(*appConfig.ConfigLocation)
	} else {
		configData, err = ioutil.ReadFile(defaultConfigFileName)
	}

	return configData, err
}
