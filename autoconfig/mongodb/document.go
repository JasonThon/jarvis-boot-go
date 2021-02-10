package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type Document interface {
	CollectionName() string
	ObjectId() primitive.ObjectID
	Init()
}
