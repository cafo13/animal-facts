package handler

import (
	"math/rand"

	"github.com/pkg/errors"

	"github.com/cafo13/animal-facts/pkg/repository"
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

func (f *FactsHandler) mapFactToHandler(fact *repository.Fact) *Fact {
	return &Fact{
		Fact:   fact.Fact,
		Source: fact.Source,
	}
}

func (f *FactsHandler) Get(id int) (*Fact, error) {
	repositoryFact, err := f.factsRepository.ReadOne(id)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get fact by ID %v", id)
	}

	return f.mapFactToHandler(repositoryFact), nil
}

func (f *FactsHandler) GetRandomApproved() (*Fact, error) {
	idsOfApprovedFacts, err := f.factsRepository.ReadManyIDs(func(fact *repository.Fact) bool {
		if fact.Approved {
			return true
		}
		return false
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not get IDs of approved facts")
	}

	if len(idsOfApprovedFacts) == 0 {
		return nil, errors.New("no approved facts found")
	}

	randomFactId := idsOfApprovedFacts[rand.Intn(len(idsOfApprovedFacts))]

	randomFact, err := f.factsRepository.ReadOne(randomFactId)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get random fact by ID %v", randomFactId)
	}

	return f.mapFactToHandler(randomFact), nil
}
