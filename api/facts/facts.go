package facts

import (
	"fmt"
	"math/rand"

	"github.com/cafo13/animal-facts/api/database"
	"github.com/cafo13/animal-facts/api/types"

	log "github.com/sirupsen/logrus"
)

type FactHandler interface {
	AddFact(fact *types.Fact) error
	DeleteFact(id string) error
	FactExists(id string) bool
	GetFactById(id string) (*types.Fact, error)
	GetRandomFact() (*types.Fact, error)
	UpdateFact(id string, dact *types.Fact) error
}

type FactDataHandler struct {
	Handler database.DatabaseHandler
}

func NewFactHandler(databaseHandler database.DatabaseHandler) FactHandler {
	return FactDataHandler{Handler: databaseHandler}
}

func (dh FactDataHandler) AddFact(fact *types.Fact) error {
	err := dh.Handler.AddItem(fact)
	if err != nil {
		return err
	}
	return nil
}

func (dh FactDataHandler) FactExists(id string) bool {
	return dh.Handler.ItemExists(id)
}

func (dh FactDataHandler) DeleteFact(id string) error {
	err := dh.Handler.DeleteItem(id)
	if err != nil {
		return err
	}
	return nil
}

func (dh FactDataHandler) UpdateFact(id string, fact *types.Fact) error {
	err := dh.Handler.UpdateItem(id, fact)
	if err != nil {
		return err
	}
	return nil
}

func (dh FactDataHandler) GetFactById(id string) (*types.Fact, error) {
	item, err := dh.Handler.GetItem(id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (dh FactDataHandler) GetRandomFact() (*types.Fact, error) {
	itemCount, err := dh.Handler.GetItemCount()
	if err != nil {
		return nil, err
	}

	if itemCount < 2 {
		return nil, fmt.Errorf("collection needs to have at least two documents to get random item, collection has %d", itemCount)
	}

	log.Infof("Total available facts in DB: %d", itemCount)

	randomId := rand.Int63n(itemCount) + 1

	log.Infof("Chosen random fact with id: %d", randomId)

	item, err := dh.Handler.GetItem(fmt.Sprintf("%d", randomId))
	if err != nil {
		return nil, err
	}
	return item, nil
}
