package handler

type Database struct {
	Db string
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
	return Fact{Id: id, Text: "All animals are awesome!", Category: "general", Source: "https://animalfacts.app/sources", Image: ""}, nil
}

func (db *Database) GetRandomFact() (Fact, error) {
	return Fact{Id: "123", Text: "All animals are awesome!", Category: "random", Source: "https://animalfacts.app/sources", Image: ""}, nil
}
