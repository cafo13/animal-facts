package handler

import (
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/cafo13/animal-facts/pkg/repository"
)

var (
	ErrNotFound = errors.New("fact not found")
)

type Fact struct {
	ID       primitive.ObjectID `json:"id"`
	Fact     string             `json:"fact"`
	Source   string             `json:"source"`
	Approved bool               `json:"approved"`
}

type FactsHandler struct {
	factsRepository repository.FactsRepository
}

func NewFactsHandler(factsRepository repository.FactsRepository) *FactsHandler {
	return &FactsHandler{factsRepository}
}

func (f *FactsHandler) mapFactToHandler(fact *repository.Fact) *Fact {
	return &Fact{
		ID:     fact.ID,
		Fact:   fact.Fact,
		Source: fact.Source,
	}
}

func (f *FactsHandler) Create(fact *Fact) error {
	factToCreate := &repository.Fact{
		ID:        fact.ID,
		Fact:      fact.Fact,
		Source:    fact.Source,
		Approved:  fact.Approved,
		CreatedAt: time.Now(),
		CreatedBy: "user.name", // TODO set user name
		UpdatedAt: time.Now(),
		UpdatedBy: "user.name", // TODO set user name
	}

	err := f.factsRepository.Create(factToCreate)
	if err != nil {
		return errors.Wrapf(err, "failed to create fact")
	}

	return nil
}

func (f *FactsHandler) Update(fact *Fact) error {
	return errors.New("not implemented")
}

func (f *FactsHandler) Delete(id primitive.ObjectID) error {
	return f.factsRepository.Delete(id)
}
