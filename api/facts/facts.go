package facts

import (
	"math/rand"

	"github.com/cafo13/animal-facts/api/database"
	"github.com/cafo13/animal-facts/api/types"
)

type DataHandler struct {
	Handler database.DatabaseHandler
}

type FactHandler interface {
	GetFactById(id string) (*types.Fact, error)
	GetRandomFact() (*types.Fact, error)
}

func NewFactHandler(databaseHandler database.DatabaseHandler) (FactHandler, error) {
	return DataHandler{Handler: databaseHandler}, nil
}

func (dh DataHandler) GetFactById(id string) (*types.Fact, error) {
	item, err := dh.Handler.GetItem(id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (dh DataHandler) GetRandomFact() (*types.Fact, error) {
	itemCount, err := dh.Handler.GetItemCount()
	if err != nil {
		return nil, err
	}

	randomId := rand.Intn(itemCount-1) + 1

	item, err := dh.Handler.GetItem(string(rune(randomId)))
	if err != nil {
		return nil, err
	}
	return &item, nil
}
