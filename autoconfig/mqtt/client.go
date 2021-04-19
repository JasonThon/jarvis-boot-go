package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	config2 "thingworks.net/thingworks/jarvis-boot/autoconfig/config"
	"time"
)

var client Client

type Client struct {
	client      mqtt.Client
	defaultConf config2.AppConfig
}

func (cli Client) IsValid() bool {
	return cli.client != nil
}

func (cli Client) Conn() (mqtt.Token, error) {
	if cli.client != nil {
		if token := cli.client.Connect(); token.Wait() && token.Error() != nil {
			logrus.Errorf("Error happens when connection mqtt client %v", token.Error())

			return token, token.Error()
		}

		logrus.Infof("Mqtt client has connected to broker [%v] successfully", cli.defaultConf.Mqtt.GetBroker())
	}

	return nil, nil
}

func (cli Client) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	if cli.IsValid() {
		// if client is not connected, we will try to connect it
		if !cli.client.IsConnected() || !cli.client.IsConnectionOpen() {
			// Here, we will retry to connect mqtt as configured times
			for i := 0; i < cli.defaultConf.Mqtt.GetRetry(); i++ {
				if _, err := cli.Conn(); err == nil {
					return cli.client.Publish(topic, qos, retained, payload)
				}
			}
		}

		return cli.client.Publish(topic, qos, retained, payload)
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
		opts := mqtt.NewClientOptions().AddBroker(defaultConfig.Mqtt.GetBroker()).SetClientID(defaultConfig.Mqtt.ClientId)
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
