package facts

import (
	"github.com/cafo13/animal-facts/api/database"
	"github.com/cafo13/animal-facts/api/models"

	"github.com/pkg/errors"
)

type FactHandler interface {
	CreateFact(fact *models.Fact) error
	DatabaseFactToModelsFactWithID(dbFact *database.Fact) *models.FactWithID
	DeleteFact(id uint) error
	GetFactById(id uint) (*models.FactWithID, error)
	GetRandomFact() (*models.FactWithID, error)
	ModelsFactToDatabaseFact(fact *models.Fact) *database.Fact
	UpdateFact(id uint, updatedFact *models.Fact) (*models.FactWithID, error)
}

type FactDataHandler struct {
	Handler database.DatabaseHandler
}

func NewFactHandler(databaseHandler database.DatabaseHandler) FactHandler {
	return FactDataHandler{Handler: databaseHandler}
}

func (fdh FactDataHandler) CreateFact(fact *models.Fact) error {
	dbFact := fdh.ModelsFactToDatabaseFact(fact)
	err := dbFact.Create()
	if err != nil {
		return err
	}

	return nil
}

func (fdh FactDataHandler) DatabaseFactToModelsFactWithID(dbFact *database.Fact) *models.FactWithID {
	fact := &models.FactWithID{}
	fact.ID = dbFact.ID
	fact.Text = dbFact.Text
	fact.Category = dbFact.Category
	fact.Source = dbFact.Source

	return fact
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

func (fdh FactDataHandler) GetFactById(id uint) (*models.FactWithID, error) {
	dbFact := &database.Fact{}
	dbFact.ID = id
	dbFact.Database = *fdh.Handler.GetDatabase()
	err := dbFact.Read()
	if err != nil {
		return nil, err
	}

	return fdh.DatabaseFactToModelsFactWithID(dbFact), nil
}

func (fdh FactDataHandler) GetRandomFact() (*models.FactWithID, error) {
	dbFact := &database.Fact{}
	dbFact.Database = *fdh.Handler.GetDatabase()
	factId, err := dbFact.GetRandomFactId()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get random fact id")
	}

	dbFact.ID = uint(factId)
	err = dbFact.Read()
	if err != nil {
		return nil, err
	}

	return fdh.DatabaseFactToModelsFactWithID(dbFact), nil
}

func (fdh FactDataHandler) ModelsFactToDatabaseFact(fact *models.Fact) *database.Fact {
	dbFact := &database.Fact{}
	dbFact.Database = *fdh.Handler.GetDatabase()
	dbFact.Text = fact.Text
	dbFact.Category = fact.Category
	dbFact.Source = fact.Source

	return dbFact
}

func (fdh FactDataHandler) UpdateFact(id uint, updatedFact *models.Fact) (*models.FactWithID, error) {
	dbFact := fdh.ModelsFactToDatabaseFact(updatedFact)
	dbFact.ID = id
	err := dbFact.Update()
	if err != nil {
		return nil, err
	}

	fact, err := fdh.GetFactById(dbFact.ID)
	if err != nil {
		return nil, err
	}

	return fact, nil
}
