package mongodb

import (
	config2 "github.com/thingworks/common/autoconfig/config"
)

func InitMongoTemplate(config config2.MongoConfig) error {
	err, conn := NewConnector(config)

	if err != nil {
		return err
	}

	template = *conn.getMongoTemplate(config.DataBase)

	return nil
}
