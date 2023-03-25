package command

import (
	"context"

	"github.com/cafo13/animal-facts/backend/common/decorator"
	"github.com/cafo13/animal-facts/backend/common/logs"
	"github.com/cafo13/animal-facts/backend/facts/domain/fact"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type DeleteFact struct {
	FactUUID uuid.UUID
}

type DeleteFactHandler decorator.CommandHandler[DeleteFact]

type deleteFactHandler struct {
	repo fact.Repository
}

func NewDeleteFactHandler(
	repo fact.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) decorator.CommandHandler[DeleteFact] {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[DeleteFact](
		deleteFactHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h deleteFactHandler) Handle(ctx context.Context, cmd DeleteFact) (err error) {
	defer func() {
		logs.LogCommandExecution("DeleteFactHandler", cmd, err)
	}()

	return h.repo.DeleteFact(ctx, cmd.FactUUID)
}
