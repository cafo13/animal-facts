package handler

import (
	"github.com/cafo13/animal-facts/pkg/repository"
	"github.com/pkg/errors"
)

type Fact struct {
	Fact   string `bson:"fact"`
	Source string `bson:"source"`
}

type FactsHandler struct {
	factsRepository repository.FactsRepository
}

func NewFactsHandler(factsRepository repository.FactsRepository) *FactsHandler {
	return &FactsHandler{factsRepository}
}

func (f FactsHandler) mapFactToHandler(fact *repository.Fact) *Fact {
	return &Fact{
		Fact:   fact.Fact,
		Source: fact.Source,
	}
}

func (f FactsHandler) Get(id int) (*Fact, error) {
	repositoryFact, err := f.factsRepository.Get(id)
	if err != nil {
		return nil, err
	}

	return f.mapFactToHandler(repositoryFact), nil
}

func (f FactsHandler) GetRandomApproved() (*Fact, error) {
	return nil, errors.New("not implemented")
}
