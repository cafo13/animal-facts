package database

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cafo13/animal-facts/api/types"
)

type Database struct {
	MongoDatabaseUri  string
	MongoDatabaseName string
}

type DatabaseHandler interface {
	GetItem(id string) (*types.Fact, error)
	GetItemCount() (int64, error)
}

func NewDatabaseHandler(mongoDatabaseUri string, mongoDatabaseName string) DatabaseHandler {
	return Database{MongoDatabaseUri: mongoDatabaseUri, MongoDatabaseName: mongoDatabaseName}
}

func (db Database) GetItem(id string) (*types.Fact, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(db.MongoDatabaseUri))
	if err != nil {
		log.Fatal("Failed to initialize mongo db client", err)
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to mongo db", err)
		return nil, err
	}
	defer client.Disconnect(ctx)

	database := client.Database(db.MongoDatabaseName)
	collection := database.Collection("animalfacts")

	result := collection.FindOne(context.Background(), bson.M{"Id": id})
	fact := &types.Fact{}
	result.Decode(fact)

	if fact.Id == "" {
		return nil, fmt.Errorf("fact with id %s not found", id)
	}

	return fact, nil
}

func (db Database) GetItemCount() (int64, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(db.MongoDatabaseUri))
	if err != nil {
		log.Fatal("Failed to initialize mongo db client", err)
		return -1, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to mongo db", err)
		return -1, err
	}
	defer client.Disconnect(ctx)

	database := client.Database(db.MongoDatabaseName)
	collection := database.Collection("animalfacts")

	count, err := collection.EstimatedDocumentCount(ctx, nil)
	if err != nil {
		return -1, err
	}

	return count, nil
}
