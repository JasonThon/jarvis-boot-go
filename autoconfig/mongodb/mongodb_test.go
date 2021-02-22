package mongodb

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/thingworks/common/autoconfig/config"
	"github.com/thingworks/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func getDatabase() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	clientOptions = clientOptions.SetAuth(options.Credential{
		Username: "jarvis",
		Password: "gpCjkZCCMZeqbZzjdJ2ELrF#",
	})

	client, err := mongo.Connect(context, clientOptions)

	if err != nil {
		return nil
	}

	return client.Database("jarvis", options.Database())
}

type Properties struct {
	CreatedBy string `bson:"createdBy" json:"createdBy"`
}

type testDocument struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string
	Time       time.Time
	Properties Properties
}

func (t *testDocument) Init() {
}

func (t *testDocument) CollectionName() string {
	return "test_go_mongo"
}

func (t *testDocument) TypeAlias() string {
	return ""
}

func (t *testDocument) ObjectId() primitive.ObjectID {
	return t.Id
}

func TestNewConnector(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Exception happens when testing mongo saving: %v", r)
		}
	}()

	err, connector := NewConnector(config.MongoConfig{
		Host:     "localhost",
		Port:     "27017",
		Username: "jarvis",
		Password: "gpCjkZCCMZeqbZzjdJ2ELrF#",
		DataBase: "jarvis",
	})

	assert.Nil(t, err)

	mongoTemplate := connector.getMongoTemplate("test_go_mongo")
	saveDoc(mongoTemplate)

	assert.NotNil(t, mongoTemplate.FindOne(bson.D{}, &testDocument{}, options.FindOne()))

	assert.Nil(t, mongoTemplate.DeleteOne(bson.D{{"Name", "Test"}}, "test_go_mongo"))
}

func TestMongoTemplate_Save(t *testing.T) {
	var testTemplate = NewMongoTemplate(getDatabase())
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Exception happens when testing mongo saving: %v", r)
		}
	}()

	saveDoc(testTemplate)
}

func saveDoc(testTemplate *MongoTemplate) {
	now, _ := utils.Parse(utils.NowString())

	testTemplate.Save(&testDocument{
		Name:       "Test",
		Time:       now,
		Properties: Properties{CreatedBy: "Me"},
	})
}

func TestMongoTemplate_FindAll(t *testing.T) {
	template := NewMongoTemplate(getDatabase())
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Eception happens when finding: %v", r)
		}
	}()

	documents := template.FindAll(
		bson.D{{"name", "Test"}}, &testDocument{},
		options.Find().SetProjection(bson.D{{"name", 1}}))

	assert.Greater(t, len(documents), 0)
	for _, document := range documents {
		assert.Equal(t, document.CollectionName(), "test_go_mongo")
		log.Info(document.(*testDocument).Id)
	}
}

func TestMongoTemplate_DeleteOne(t *testing.T) {
	template := NewMongoTemplate(getDatabase())
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Eception happens when finding: %v", r)
		}
	}()

	err := template.DeleteAll(
		bson.D{
			{"resource.resourceId", bson.D{
				{"$in", []string{"6275477A-A0BE-495E-906A-CC9B82539EC2"}}},
			},
			{"resource.subResourceId", bson.D{
				{"$in", []string{"d3d6b572"}}},
			},
		},
		"header.entries.changeLog")

	if err != nil {
		panic(err)
	}
}
