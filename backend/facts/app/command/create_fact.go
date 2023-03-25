package command

import (
	"context"

	"github.com/cafo13/animal-facts/backend/common/decorator"
	"github.com/cafo13/animal-facts/backend/common/logs"
	"github.com/cafo13/animal-facts/backend/facts/domain/fact"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CreateFact struct {
	FactUUID   uuid.UUID
	FactText   string
	FactSource string
}

type CreateFactHandler decorator.CommandHandler[CreateFact]

type createFactHandler struct {
	repo fact.Repository
}

func NewCreateFactHandler(
	repo fact.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) decorator.CommandHandler[CreateFact] {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[CreateFact](
		createFactHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h createFactHandler) Handle(ctx context.Context, cmd CreateFact) (err error) {
	defer func() {
		logs.LogCommandExecution("CreateFactHandler", cmd, err)
	}()

	f, err := fact.NewFact(cmd.FactUUID, cmd.FactText, cmd.FactSource)
	if err != nil {
		return err
	}

	return h.repo.AddFact(ctx, f)
}
