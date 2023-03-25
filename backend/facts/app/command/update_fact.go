package command

import (
	"context"

	"github.com/cafo13/animal-facts/backend/common/decorator"
	"github.com/cafo13/animal-facts/backend/common/logs"
	"github.com/cafo13/animal-facts/backend/facts/domain/fact"
	"github.com/sirupsen/logrus"
)

type UpdateFact struct {
	FactUUID string
	NewText  string
	NewSouce string
}

type UpdateFactHandler decorator.CommandHandler[UpdateFact]

type updateFactHandler struct {
	repo fact.Repository
}

func NewUpdateFactHandler(
	repo fact.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) decorator.CommandHandler[UpdateFact] {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[UpdateFact](
		updateFactHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h updateFactHandler) Handle(ctx context.Context, cmd UpdateFact) (err error) {
	defer func() {
		logs.LogCommandExecution("UpdateFactHandler", cmd, err)
	}()

	return h.repo.UpdateFact(
		ctx,
		cmd.FactUUID,
		func(ctx context.Context, f *fact.Fact) (*fact.Fact, error) {
			if err := f.UpdateText(cmd.NewText); err != nil {
				return nil, err
			}

			if err := f.UpdateSource(cmd.NewSouce); err != nil {
				return nil, err
			}

			return f, nil
		},
	)
}
