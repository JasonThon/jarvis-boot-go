package mongodb

import (
	context2 "context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	context        = context2.TODO()
	collectionOpts = options.Collection()
	template       MongoTemplate
)

type MongoTemplate struct {
	database *mongo.Database
}

func NewMongoTemplate(database *mongo.Database) *MongoTemplate {
	return &MongoTemplate{
		database: database,
	}
}

func (template *MongoTemplate) FindAll(filter interface{}, document Document, opts ...*options.FindOptions) []Document {
	collection := template.collection(document.CollectionName())
	cursor, err := collection.Find(context, filter, opts...)

	if err != nil {
		return nil
	}

	defer closeCursor(cursor)

	docs := make([]Document, 0)

	for cursor.Next(context) {
		err := cursor.Decode(document)

		if err != nil {
			return nil
		}

		docs = append(docs, document)
	}

	return docs

}

func (template *MongoTemplate) FindAllAndAssign(filter interface{}, collectionName string, document []Document,
	opts ...*options.FindOptions) error {
	collection := template.collection(collectionName)
	cursor, err := collection.Find(context, filter, opts...)

	if err != nil {
		return err
	}

	if cursor != nil && checkCollection(collectionName, document) {
		err = cursor.All(context, document)

		if err != nil {
			return err
		}
	}

	return nil
}

func checkCollection(name string, documents []Document) bool {
	for _, doc := range documents {
		if doc.CollectionName() != name {
			return false
		}
	}

	return true
}

func closeCursor(cursor *mongo.Cursor) {
	err := cursor.Close(context)

	if err != nil {
		log.WithFields(log.Fields{
			"context": context,
		}).Errorf("Exception happens when close cursor: %v", err)
	}
}

func (template *MongoTemplate) FindOne(filter interface{}, result Document, opts ...*options.FindOneOptions) Document {
	collection := template.collection(result.CollectionName())
	singleResult := collection.FindOne(context, filter, opts...)
	err := singleResult.Decode(result)

	if err != nil {
		log.WithFields(log.Fields{
			"query": filter,
		}).Errorf("Exception when query single result. The detail is: %v", err)
		return nil
	}

	return result
}

func (template *MongoTemplate) DeleteOne(filter interface{}, collectionName string) error {
	collection := template.collection(collectionName)
	_, err := collection.DeleteOne(context, filter, options.Delete())

	if err != nil {
		log.WithFields(log.Fields{
			"query": filter,
		}).Errorf("Exception when delete single document. The detail is: %v", err)
	}

	return err
}

func (template *MongoTemplate) DeleteAll(filter interface{}, collectionName string) error {
	collection := template.collection(collectionName)
	_, err := collection.DeleteMany(context, filter, options.Delete())

	if err != nil {
		log.WithFields(log.Fields{
			"query": filter,
		}).Errorf("Exception when delete all documents. The detail is: %v", err)
	}

	return err
}

func (template *MongoTemplate) UpdateOne(collectionName string, filter interface{}, update interface{}) error {
	collection := template.collection(collectionName)
	_, err := collection.UpdateOne(context, filter, update, options.Update())

	if err != nil {
		log.WithFields(log.Fields{
			"query":  filter,
			"update": update,
		}).Errorf("Exception when update document. The detail is: %v", err)
	}

	return err
}

func (template *MongoTemplate) UpdateMulti(collectionName string, filter interface{}, update interface{}) {
	collection := template.collection(collectionName)
	_, err := collection.UpdateMany(context, filter, update, options.Update())

	if err != nil {
		log.WithFields(log.Fields{
			"query":  filter,
			"update": update,
		}).Errorf("Exception when update all document. The detail is: %v", err)
	}
}

func (template *MongoTemplate) FindAndModify(document Document, filter interface{}, modify interface{},
	opts *options.FindOneAndUpdateOptions) Document {

	collection := template.collection(document.CollectionName())
	result := collection.FindOneAndUpdate(context, filter, modify, opts)

	err := result.Decode(document)

	if err != nil {

		log.WithFields(log.Fields{
			"query":  filter,
			"update": modify,
		}).Errorf("Exception happens. The detail is: %v", err)

		return nil
	}

	return document
}

func (template *MongoTemplate) FindAndReplace(document Document, filter interface{}, replacement interface{},
	opts *options.FindOneAndReplaceOptions) Document {
	collection := template.collection(document.CollectionName())
	result := collection.FindOneAndReplace(context, filter, replacement, opts)

	err := result.Decode(document)

	if err != nil {

		log.WithFields(log.Fields{
			"query":       filter,
			"replacement": replacement,
		}).Errorf("Exception happens. The detail is: %v", err)

		return nil
	}

	return document
}

func (template *MongoTemplate) collection(collectionName string) *mongo.Collection {
	return template.database.Collection(collectionName, collectionOpts)
}

func (template *MongoTemplate) Save(document Document) {

	if document != nil {
		collection := template.collection(document.CollectionName())
		_, err := collection.InsertOne(context, document, options.InsertOne())

		if err != nil {
			log.WithFields(log.Fields{
				"collection": document.CollectionName(),
			}).Errorf("Document save failed: %v", err)
		}
	}
}

func DefaultMongoTemplate() *MongoTemplate {
	return &template
}
