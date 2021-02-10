package mongodb

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	config2 "thingworks/common/autoconfig/config"
	"thingworks/common/utils/strings2"
)

type connectionError struct{ err error }

func (c connectionError) Error() string {
	return fmt.Sprintf("Exception happens when connect MongoDB: %v", c.err)
}

type MongoDBConnector struct {
	client *mongo.Client
}

func NewConnector(mongoConf config2.MongoConfig) (error, *MongoDBConnector) {

	clientOpts := getOptions(mongoConf)

	client, err := mongo.Connect(context, clientOpts)

	logrus.Infof("Mongo Connection build success. The uri is: %s", mongoConf.GetUri())

	return err, &MongoDBConnector{
		client: client,
	}
}

func (conn *MongoDBConnector) getMongoTemplate(databaseName string) *MongoTemplate {
	database := conn.client.Database(databaseName, options.Database())
	logrus.Infof("Mongo Database [%s] is connected", databaseName)
	return NewMongoTemplate(database)
}

func  getOptions(mongoConf config2.MongoConfig) *options.ClientOptions {
	clientOptions := &options.ClientOptions{}
	var credential options.Credential

	if strings2.IsNotBlank(mongoConf.Username) && strings2.IsNotBlank(mongoConf.Password) {

		credential = options.Credential{
			Username: mongoConf.Username,
			Password: mongoConf.Password,
		}

		clientOptions = clientOptions.SetAuth(credential)
	}

	return clientOptions.
		ApplyURI(mongoConf.GetUri()).
		SetConnectTimeout(mongoConf.Timeout)
}
