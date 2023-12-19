package handler

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"

	"github.com/cafo13/animal-facts/pkg/repository"
)

var (
	ErrNotFound = errors.New("fact not found")
)

type Fact struct {
	ID     string `bson:"id" json:"id"`
	Fact   string `bson:"fact" json:"fact"`
	Source string `bson:"source" json:"source"`
}

type FactsHandler struct {
	factsRepository repository.FactsRepository
}

func NewFactsHandler(factsRepository repository.FactsRepository) *FactsHandler {
	return &FactsHandler{factsRepository}
}

func (f *FactsHandler) mapFactToHandler(fact *repository.Fact) *Fact {
	return &Fact{
		ID:     fact.ID.Hex(),
		Fact:   fact.Fact,
		Source: fact.Source,
	}
}

func (f *FactsHandler) Get(id primitive.ObjectID) (*Fact, error) {
	repositoryFact, err := f.factsRepository.ReadOne(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, errors.Wrapf(err, "could not get fact by ID %v", id)
	}

	return f.mapFactToHandler(repositoryFact), nil
}

func (f *FactsHandler) GetRandomApproved() (*Fact, error) {
	idsOfApprovedFacts, err := f.factsRepository.ReadManyIDs(func(fact *repository.Fact) bool {
		return true
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

func (f *FactsHandler) GetFactsCount() (int, error) {
	factsCount, err := f.factsRepository.Count()
	if err != nil {
		return 0, errors.Wrapf(err, "could not get facts count")
	}

	return factsCount, nil
}
