package repository

import (
	"context"
	"github.com/pkg/errors"
	"time"

	"github.com/neko-neko/echo-logrus/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Fact struct {
	ID        int       `bson:"_id"`
	Fact      string    `bson:"fact"`
	Source    string    `bson:"source"`
	CreatedAt time.Time `bson:"created_at"`
	CreatedBy string    `bson:"created_by"`
	UpdatedAt time.Time `bson:"updated_at"`
	UpdatedBy string    `bson:"updated_by"`
}

type FactsRepository interface {
	GetFact(id int) (*Fact, error)
	GetRandomFact() (*Fact, error)
	GetFactCount() (uint64, error)
}

type MongoDBFactsRepository struct {
	mongoDbClient *mongo.Client
}

func NewMongoDBFactsRepository(mongoDbUri string) (FactsRepository, error) {
	opts := options.Client().ApplyURI(mongoDbUri).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Logger().WithError(err).Fatal("failed to disconnect from mongo db")
		}
	}()

	if err := client.Database("animal-facts").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, errors.Wrap(err, "failed to ping mongo db")
	}
	log.Logger().Info("connected to mongo db")

	return &MongoDBFactsRepository{client}, nil
}

func (m *MongoDBFactsRepository) GetFact(id int) (*Fact, error) {
	filter := bson.D{{"_id", id}}
	var result Fact
	err := m.factsCollection().FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *MongoDBFactsRepository) GetRandomFact() (*Fact, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBFactsRepository) GetFactCount() (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBFactsRepository) factsCollection() *mongo.Collection {
	return m.mongoDbClient.Database("animal-facts").Collection("facts")
}
