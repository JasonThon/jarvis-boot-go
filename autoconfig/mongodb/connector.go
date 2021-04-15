package mongodb

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	config2 "thingworks.net/thingworks/jarvis-boot/autoconfig/config"
	"thingworks.net/thingworks/jarvis-boot/utils/strings2"
)

type Connector struct {
	client *mongo.Client
}

func NewConnector(mongoConf config2.MongoConfig) (error, *Connector) {
	if mongoConf.IsValidConf() {
		clientOpts := getOptions(mongoConf)

		client, err := mongo.Connect(context, clientOpts)

		logrus.Infof("Mongo Connection build success. The uri is: %s", mongoConf.GetUri())

		return err, &Connector{
			client: client,
		}
	}

	return nil, nil
}

func (conn *Connector) getMongoTemplate(databaseName string) *MongoTemplate {
	database := conn.client.Database(databaseName, options.Database())
	logrus.Infof("Mongo Database [%s] is connected", databaseName)
	return NewMongoTemplate(database)
}

func getOptions(mongoConf config2.MongoConfig) *options.ClientOptions {
	clientOptions := &options.ClientOptions{}
	var credential options.Credential

	if strings2.IsNotBlank(mongoConf.Username) && strings2.IsNotBlank(mongoConf.Password) {

		credential = options.Credential{
			Username:   mongoConf.Username,
			Password:   mongoConf.Password,
			AuthSource: mongoConf.DataBase,
		}
	}

	return clientOptions.
		ApplyURI(mongoConf.GetUri()).
		SetAuth(credential)
}
