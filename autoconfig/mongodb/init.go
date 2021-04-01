package mongodb

import (
	config2 "thingworks.net/thingworks/common/autoconfig/config"
)

func InitMongoTemplate(config config2.MongoConfig) error {
	err, conn := NewConnector(config)

	if err != nil {
		return err
	}

	if conn != nil {
		template = *conn.getMongoTemplate(config.DataBase)
	}

	return nil
}
