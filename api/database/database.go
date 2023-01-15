package database

import (
	"context"
	"fmt"

	"github.com/cafo13/animal-facts/api/types"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseHandler interface {
	AddItem(fact *types.Fact) error
	CloseConnection() error
	DeleteItem(id string) error
	GetItem(id string) (*types.Fact, error)
	GetItemCount() (int64, error)
	ItemExists(id string) (bool, error)
	UpdateItem(id string, fact *interface{}) (*types.Fact, error)
}

type Database struct {
	Client     *mongo.Client
	Context    *context.Context
	Collection *mongo.Collection
}

func NewDatabaseHandler(mongoDatabaseUri string, mongoDatabaseName string, mongoCollectionName string) (DatabaseHandler, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDatabaseUri))
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize mongo db client")
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongo db")
	}

	database := client.Database(mongoDatabaseName)
	collection := database.Collection(mongoCollectionName)

	return Database{Client: client, Collection: collection, Context: &ctx}, nil
}

func (db Database) AddItem(fact *types.Fact) error {
	_, err := db.Collection.InsertOne(context.Background(), fact)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to insert new item to database, item: %+v", fact))
	}

	return nil
}

func (db Database) CloseConnection() error {
	err := db.Client.Disconnect(*db.Context)
	if err != nil {
		return errors.Wrap(err, "failed to disconnect from the database")
	}

	return nil
}

func (db Database) DeleteItem(id string) error {
	_, err := db.Collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to delete item from database, item with id %s", id))
	}

	return nil
}

func (db Database) GetItem(id string) (*types.Fact, error) {
	result := db.Collection.FindOne(context.Background(), bson.M{"id": id})
	fact := &types.Fact{}
	result.Decode(fact)

	if fact.Id == "" {
		return nil, fmt.Errorf("fact with id %s not found", id)
	}

	return fact, nil
}

func (db Database) GetItemCount() (int64, error) {
	count, err := db.Collection.EstimatedDocumentCount(*db.Context, nil)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (db Database) ItemExists(id string) (bool, error) {
	result := db.Collection.FindOne(context.Background(), bson.M{"id": id})
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (db Database) UpdateItem(id string, fact *interface{}) (*types.Fact, error) {
	_, err := db.Collection.UpdateOne(context.Background(), bson.M{"id": id}, bson.M{"$set": fact})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to update item %s in database, item: %+v", id, fact))
	}

	updatedFact, err := db.GetItem(id)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get updated item %s from database", id))
	}

	return updatedFact, nil
}
