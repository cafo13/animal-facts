package repository

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/neko-neko/echo-logrus/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNotFound = errors.New("fact not found")
)

type Fact struct {
	ID        primitive.ObjectID `bson:"_id"`
	Fact      string             `bson:"fact"`
	Source    string             `bson:"source"`
	Approved  bool               `bson:"approved"`
	CreatedAt time.Time          `bson:"created_at"`
	CreatedBy string             `bson:"created_by"`
	UpdatedAt time.Time          `bson:"updated_at"`
	UpdatedBy string             `bson:"updated_by"`
}

type FactsRepository interface {
	Create(fact *Fact) error
	ReadOne(id primitive.ObjectID) (*Fact, error)
	ReadMany(filterFunc func(fact *Fact) bool) ([]*Fact, error)
	ReadManyIDs(filterFunc func(fact *Fact) bool) ([]primitive.ObjectID, error)
	Update(id primitive.ObjectID, updateFunc func(fact *Fact) error) error
	Delete(id primitive.ObjectID) error
	Count() (int, error)
	Close() error
}

type MongoDBFactsRepository struct {
	mongoDbClient *mongo.Client
	databaseName  string
}

func NewMongoDBFactsRepository(mongoDbUri string) (FactsRepository, error) {
	opts := options.Client().ApplyURI(mongoDbUri).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	databaseName := "animal-facts"

	databaseNameFromEnv, exists := os.LookupEnv("MONGODB_DATABASE_NAME")
	if exists {
		databaseName = databaseNameFromEnv
	}

	log.Logger().Info("using database: " + databaseName)
	if err := client.Database(databaseName).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, errors.Wrap(err, "failed to ping mongo db")
	}
	log.Logger().Info("connected to mongo db")

	return &MongoDBFactsRepository{client, databaseName}, nil
}

func (m *MongoDBFactsRepository) factsCollection() *mongo.Collection {
	return m.mongoDbClient.Database(m.databaseName).Collection("facts")
}

func (m *MongoDBFactsRepository) Create(fact *Fact) error {
	_, err := m.factsCollection().InsertOne(context.TODO(), fact)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDBFactsRepository) ReadOne(id primitive.ObjectID) (*Fact, error) {
	filter := bson.D{{"_id", id}, {"approved", true}}
	var result Fact
	err := m.factsCollection().FindOne(context.TODO(), filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *MongoDBFactsRepository) ReadMany(filterFunc func(fact *Fact) bool) ([]*Fact, error) {
	return nil, errors.New("not implemented")
}

func (m *MongoDBFactsRepository) ReadManyIDs(filterFunc func(fact *Fact) bool) ([]primitive.ObjectID, error) {
	filter := bson.D{{"approved", true}}
	var facts []Fact
	cursor, err := m.factsCollection().Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &facts); err != nil {
		return nil, err
	}

	var result []primitive.ObjectID
	for _, fact := range facts {
		if filterFunc(&fact) {
			result = append(result, fact.ID)
		}
	}

	return result, nil
}

func (m *MongoDBFactsRepository) Update(id primitive.ObjectID, updateFunc func(fact *Fact) error) error {
	return errors.New("not implemented")
}

func (m *MongoDBFactsRepository) Delete(id primitive.ObjectID) error {
	return errors.New("not implemented")
}

func (m *MongoDBFactsRepository) Count() (int, error) {
	filter := bson.D{{"approved", true}}
	var facts []Fact
	cursor, err := m.factsCollection().Find(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	if err = cursor.All(context.TODO(), &facts); err != nil {
		return 0, err
	}

	return len(facts), nil
}

func (m *MongoDBFactsRepository) Close() error {
	if err := m.mongoDbClient.Disconnect(context.TODO()); err != nil {
		log.Logger().WithError(err).Fatal("failed to disconnect from mongo db")
		return err
	}

	return nil
}

type MockFactsRepository struct {
	facts                 map[primitive.ObjectID]*Fact
	errorAllFunctionCalls bool
}

func NewMockFactsRepository(facts map[primitive.ObjectID]*Fact, errorAllFunctionCalls bool) FactsRepository {
	return &MockFactsRepository{facts, errorAllFunctionCalls}
}

func (m *MockFactsRepository) Create(fact *Fact) error {
	if m.errorAllFunctionCalls {
		return errors.New("error at creating fact")
	}

	return nil
}

func (m *MockFactsRepository) ReadOne(id primitive.ObjectID) (*Fact, error) {
	if m.errorAllFunctionCalls {
		return nil, errors.New("error at getting fact")
	}

	if fact, exists := m.facts[id]; exists {
		return fact, nil
	}

	return nil, errors.New("fact not found")
}

func (m *MockFactsRepository) ReadMany(filterFunc func(fact *Fact) bool) ([]*Fact, error) {
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

func (m *MockFactsRepository) ReadManyIDs(filterFunc func(fact *Fact) bool) ([]primitive.ObjectID, error) {
	if m.errorAllFunctionCalls {
		return nil, errors.New("error at getting fact IDs")
	}

	var matchingFactIDs []primitive.ObjectID
	for id, fact := range m.facts {
		if filterFunc(fact) {
			matchingFactIDs = append(matchingFactIDs, id)
		}
	}

	return matchingFactIDs, nil
}

func (m *MockFactsRepository) Update(id primitive.ObjectID, updateFunc func(fact *Fact) error) error {
	if m.errorAllFunctionCalls {
		return errors.New("error at updating fact")
	}

	return nil
}

func (m *MockFactsRepository) Delete(id primitive.ObjectID) error {
	if m.errorAllFunctionCalls {
		return errors.New("error at deleting fact")
	}

	return nil
}

func (m *MockFactsRepository) Count() (int, error) {
	if m.errorAllFunctionCalls {
		return 0, errors.New("error at getting fact count")
	}

	return len(m.facts), nil
}

func (m *MockFactsRepository) Close() error {
	if m.errorAllFunctionCalls {
		return errors.New("error at getting fact count")
	}

	return nil
}
