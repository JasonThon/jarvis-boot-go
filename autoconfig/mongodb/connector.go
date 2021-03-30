package mongodb

import (
	"github.com/sirupsen/logrus"
	config2 "thingworks.net/thingworks/common/autoconfig/config"
	"thingworks.net/thingworks/common/utils/strings2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"regexp"
)

var reg = regexp.MustCompile("\\\\$\\\\{([^}]*)\\\\}")

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
