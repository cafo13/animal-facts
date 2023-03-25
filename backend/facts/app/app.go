package app

import (
	"github.com/cafo13/animal-facts/backend/facts/app/command"
	"github.com/cafo13/animal-facts/backend/facts/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateFact command.CreateFactHandler
	UpdateFact command.UpdateFactHandler
	DeleteFact command.DeleteFactHandler
}

type Queries struct {
	RandomFact query.RandomFactHandler
	FactByID   query.FactByIDHandler
}
