package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	config2 "thingworks.net/thingworks/common/autoconfig/config"
	"time"
)

var client Client

type Client struct {
	client      mqtt.Client
	defaultConf config2.AppConfig
}

func (cli Client) Client() mqtt.Client {
	return cli.client
}

func (cli Client) Conn() error {
	if cli.client != nil {
		if token := cli.client.Connect(); token.Wait() && token.Error() != nil {
			logrus.Errorf("Error happens when connection mqtt client %v", token.Error())

			return token.Error()
		}

		logrus.Infof("Mqtt client has connected to broker [%v] successfully", cli.defaultConf.Mqtt.Broker)
	}

	return nil
}

func NewClient(defaultConfig config2.AppConfig) *Client {
	return &Client{
		client:      newMqttClient(defaultConfig),
		defaultConf: defaultConfig,
	}
}

func newMqttClient(defaultConfig config2.AppConfig) mqtt.Client {
	if defaultConfig.Mqtt.IsValid() {
		opts := mqtt.NewClientOptions().AddBroker(defaultConfig.Mqtt.Broker).SetClientID(defaultConfig.Mqtt.ClientId)
		opts.SetKeepAlive(defaultConfig.Mqtt.KeepAlive * time.Second)
		opts.SetPingTimeout(defaultConfig.Mqtt.PingTimeout * time.Second)
		opts.SetConnectTimeout(defaultConfig.Mqtt.ConnectionTimeout * time.Second)
		opts.SetUsername(defaultConfig.Mqtt.Username)
		opts.SetPassword(defaultConfig.Mqtt.Password)

		cli := mqtt.NewClient(opts)

		logrus.Info("Mqtt client has been established")

		return cli
	}

	return nil
}

func Init(cli *Client) {
	client = *cli
}

func GetMqttClient() Client {
	if &client != nil {
		return client
	}

	logrus.Error("No valid mqtt client")

	return Client{}
}
