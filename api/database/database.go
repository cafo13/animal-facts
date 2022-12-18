package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDatabase(mongoDatabaseUri string, mongoDatabaseName string) (*mongo.Database, error) {
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

	return database, nil
}
