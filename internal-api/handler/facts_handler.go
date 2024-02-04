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
	err := f.factsRepository.Update(fact.ID, func(f *repository.Fact) *repository.Fact {
		if fact.Fact != f.Fact {
			f.Fact = fact.Fact
		}
		if fact.Source != f.Source {
			f.Source = fact.Source
		}
		f.UpdatedAt = time.Now()
		f.UpdatedBy = "user.name" // TODO set user name
		return f
	})
	if err != nil {
		return errors.Wrapf(err, "failed to update fact")
	}

	return nil
}

func (f *FactsHandler) Approve(factID primitive.ObjectID) error {
	err := f.factsRepository.Update(factID, func(f *repository.Fact) *repository.Fact {
		if !f.Approved {
			f.Approved = true
		}
		f.UpdatedAt = time.Now()
		f.UpdatedBy = "user.name" // TODO set user name
		return f
	})
	if err != nil {
		return errors.Wrapf(err, "failed to approve fact")
	}

	return nil
}

func (f *FactsHandler) Unapprove(factID primitive.ObjectID) error {
	err := f.factsRepository.Update(factID, func(f *repository.Fact) *repository.Fact {
		if f.Approved {
			f.Approved = false
		}
		f.UpdatedAt = time.Now()
		f.UpdatedBy = "user.name" // TODO set user name
		return f
	})
	if err != nil {
		return errors.Wrapf(err, "failed to unapprove fact")
	}

	return nil
}

func (f *FactsHandler) Delete(id primitive.ObjectID) error {
	return f.factsRepository.Delete(id)
}

func (f *FactsHandler) GetAll() ([]*Fact, error) {
	repositoryFacts, err := f.factsRepository.ReadAll()
	if err != nil {
		return nil, errors.Wrapf(err, "could not get all facts")
	}

	var facts []*Fact
	for _, repositoryFact := range repositoryFacts {
		facts = append(facts, f.mapFactToHandler(repositoryFact))
	}

	return facts, nil
}
