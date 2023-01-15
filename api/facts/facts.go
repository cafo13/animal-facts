package facts

import (
	"math/rand"

	"github.com/cafo13/animal-facts/api/database"

	log "github.com/sirupsen/logrus"
)

type FactHandler interface {
	CreateFact(fact *database.Fact) error
	DeleteFact(id uint) error
	GetFactById(id uint) (*database.Fact, error)
	GetRandomFact() (*database.Fact, error)
	UpdateFact(id uint, fact *database.Fact) (*database.Fact, error)
}

type FactDataHandler struct {
	Handler database.DatabaseHandler
}

func NewFactHandler(databaseHandler database.DatabaseHandler) FactHandler {
	return FactDataHandler{Handler: databaseHandler}
}

func (fdh FactDataHandler) CreateFact(fact *database.Fact) error {
	err := fact.Create()
	if err != nil {
		return err
	}

	return nil
}

func (fdh FactDataHandler) UpdateFact(id uint, updatedFact *database.Fact) (*database.Fact, error) {
	fact := updatedFact
	fact.ID = id
	err := fact.Update()
	if err != nil {
		return nil, err
	}

	return fact, nil
}

func (fdh FactDataHandler) DeleteFact(id uint) error {
	var fact *database.Fact
	fact.ID = id
	err := fact.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (fdh FactDataHandler) GetFactById(id uint) (*database.Fact, error) {
	var fact *database.Fact
	fact.ID = id
	err := fact.Read()
	if err != nil {
		return nil, err
	}

	return fact, nil
}

func (fdh FactDataHandler) GetRandomFact() (*database.Fact, error) {
	var fact *database.Fact

	factCount, err := fact.Count()
	if err != nil {
		return nil, err
	}

	log.Infof("total available facts in DB: %d", factCount)

	randomId := rand.Int63n(factCount) + 1

	log.Infof("chosen random fact with id: %d", randomId)

	fact.ID = uint(randomId)
	err = fact.Read()
	if err != nil {
		return nil, err
	}

	return fact, nil
}
