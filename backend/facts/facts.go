package facts

import (
	"github.com/cafo13/animal-facts/api/database"

	"github.com/pkg/errors"
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
	fact.Database = *fdh.Handler.GetDatabase()
	err := fact.Create()
	if err != nil {
		return err
	}

	return nil
}

func (fdh FactDataHandler) UpdateFact(id uint, updatedFact *database.Fact) (*database.Fact, error) {
	fact := updatedFact
	fact.ID = id
	fact.Database = *fdh.Handler.GetDatabase()
	err := fact.Update()
	if err != nil {
		return nil, err
	}

	return fact, nil
}

func (fdh FactDataHandler) DeleteFact(id uint) error {
	fact := &database.Fact{}
	fact.ID = id
	fact.Database = *fdh.Handler.GetDatabase()
	err := fact.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (fdh FactDataHandler) GetFactById(id uint) (*database.Fact, error) {
	fact := &database.Fact{}
	fact.ID = id
	fact.Database = *fdh.Handler.GetDatabase()
	err := fact.Read()
	if err != nil {
		return nil, err
	}

	return fact, nil
}

func (fdh FactDataHandler) GetRandomFact() (*database.Fact, error) {
	fact := &database.Fact{}
	fact.Database = *fdh.Handler.GetDatabase()
	factId, err := fact.GetRandomFactId()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get random fact id")
	}

	fact.ID = uint(factId)
	err = fact.Read()
	if err != nil {
		return nil, err
	}

	return fact, nil
}
