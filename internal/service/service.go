package service

import (
	"github.com/pkg/errors"

	"github.com/cafo13/animal-facts/internal/repository"
)

type Service struct {
	factsRepository repository.FactsRepository
}

func NewService(factsRepository repository.FactsRepository) *Service {
	return &Service{factsRepository}
}

func (s *Service) Start() error {
	return errors.New("not implemented")
}
