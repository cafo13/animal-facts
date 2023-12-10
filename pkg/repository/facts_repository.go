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
	Approved  bool      `bson:"approved"`
	CreatedAt time.Time `bson:"created_at"`
	CreatedBy string    `bson:"created_by"`
	UpdatedAt time.Time `bson:"updated_at"`
	UpdatedBy string    `bson:"updated_by"`
}

type FactsRepository interface {
	Create(fact *Fact) error
	ReadOne(id int) (*Fact, error)
	ReadMany(filterFunc func(fact *Fact) bool) ([]*Fact, error)
	ReadManyIDs(filterFunc func(fact *Fact) bool) ([]int, error)
	Update(id int, updateFunc func(fact *Fact) error) error
	Delete(id int) error
	Count() (int, error)
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

func (m *MongoDBFactsRepository) factsCollection() *mongo.Collection {
	return m.mongoDbClient.Database("animal-facts").Collection("facts")
}

func (m *MongoDBFactsRepository) Create(fact *Fact) error {
	return errors.New("not implemented")
}

func (m *MongoDBFactsRepository) ReadOne(id int) (*Fact, error) {
	filter := bson.D{{"_id", id}}
	var result Fact
	err := m.factsCollection().FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *MongoDBFactsRepository) ReadMany(filterFunc func(fact *Fact) bool) ([]*Fact, error) {
	return nil, errors.New("not implemented")
}

func (m *MongoDBFactsRepository) ReadManyIDs(filterFunc func(fact *Fact) bool) ([]int, error) {
	return nil, errors.New("not implemented")
}

func (m *MongoDBFactsRepository) Update(id int, updateFunc func(fact *Fact) error) error {
	return errors.New("not implemented")
}

func (m *MongoDBFactsRepository) Delete(id int) error {
	return errors.New("not implemented")
}

func (m *MongoDBFactsRepository) Count() (int, error) {
	return 1, errors.New("not implemented")
}

type MockFactsRepository struct {
	facts                 map[int]*Fact
	errorAllFunctionCalls bool
}

func NewMockFactsRepository(facts map[int]*Fact, errorAllFunctionCalls bool) FactsRepository {
	return &MockFactsRepository{facts, errorAllFunctionCalls}
}

func (m MockFactsRepository) Create(fact *Fact) error {
	if m.errorAllFunctionCalls {
		return errors.New("error at creating fact")
	}

	return nil
}

func (m MockFactsRepository) ReadOne(id int) (*Fact, error) {
	if m.errorAllFunctionCalls {
		return nil, errors.New("error at getting fact")
	}

	if fact, exists := m.facts[id]; exists {
		return fact, nil
	}

	return nil, errors.New("fact not found")
}

func (m MockFactsRepository) ReadMany(filterFunc func(fact *Fact) bool) ([]*Fact, error) {
	if m.errorAllFunctionCalls {
		return nil, errors.New("error at getting facts")
	}

	var matchingFacts []*Fact
	for _, fact := range m.facts {
		if filterFunc(fact) {
			matchingFacts = append(matchingFacts, fact)
		}
	}

	return matchingFacts, nil
}

func (m MockFactsRepository) ReadManyIDs(filterFunc func(fact *Fact) bool) ([]int, error) {
	if m.errorAllFunctionCalls {
		return nil, errors.New("error at getting fact IDs")
	}

	var matchingFactIDs []int
	for id, fact := range m.facts {
		if filterFunc(fact) {
			matchingFactIDs = append(matchingFactIDs, id)
		}
	}

	return matchingFactIDs, nil
}

func (m MockFactsRepository) Update(id int, updateFunc func(fact *Fact) error) error {
	if m.errorAllFunctionCalls {
		return errors.New("error at updating fact")
	}

	return nil
}

func (m MockFactsRepository) Delete(id int) error {
	if m.errorAllFunctionCalls {
		return errors.New("error at deleting fact")
	}

	return nil
}

func (m MockFactsRepository) Count() (int, error) {
	if m.errorAllFunctionCalls {
		return 0, errors.New("error at getting fact count")
	}

	return len(m.facts), nil
}
