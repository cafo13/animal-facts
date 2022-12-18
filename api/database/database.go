package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cafo13/animal-facts/api/types"
)

type Database struct {
	MongoDB *mongo.Database
}

type DatabaseHandler interface {
	GetItem(id string) (types.Fact, error)
	GetItemCount() (int, error)
}

func NewDatabaseHandler(mongoDatabaseUri string, mongoDatabaseName string) (DatabaseHandler, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDatabaseUri))
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

	database := client.Database(mongoDatabaseName)

	defer client.Disconnect(ctx)

	return Database{MongoDB: database}, nil
}

func (db Database) GetItem(id string) (types.Fact, error) {
	return types.Fact{Id: id, Text: "All animals are awesome!", Category: "general", Source: "https://github.com/cafo13/animal-facts/apisources", Image: ""}, nil
}

func (db Database) GetItemCount() (int, error) {
	return 111, nil
}
