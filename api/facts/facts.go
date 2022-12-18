package facts

import "go.mongodb.org/mongo-driver/mongo"

type Database struct {
	Db mongo.Database
}

type Fact struct {
	Id       string
	Text     string
	Category string
	Source   string
	Image    string
}

type FactHandler interface {
	GetFactById(id string) (Fact, error)
	GetRandomFact() (Fact, error)
}

func NewFactHandler(db Database) (FactHandler, error) {
	return &db, nil
}

func (db *Database) GetFactById(id string) (Fact, error) {
	return Fact{Id: id, Text: "All animals are awesome!", Category: "general", Source: "https://github.com/cafo13/animal-facts/apisources", Image: ""}, nil
}

func (db *Database) GetRandomFact() (Fact, error) {
	return Fact{Id: "123", Text: "All animals are awesome!", Category: "random", Source: "https://github.com/cafo13/animal-facts/apisources", Image: ""}, nil
}
