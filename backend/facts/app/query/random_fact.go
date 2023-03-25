package query

import (
	"context"

	"github.com/cafo13/animal-facts/backend/common/decorator"
	"github.com/sirupsen/logrus"
)

type RandomFact struct{}

type RandomFactHandler decorator.QueryHandler[RandomFact, Fact]

type randomFactHandler struct {
	readModel RandomFactReadModel
}

func NewRandomFactHandler(
	readModel RandomFactReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) RandomFactHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[RandomFact, Fact](
		randomFactHandler{readModel: readModel},
		logger,
		metricsClient,
	)
}

type RandomFactReadModel interface {
	FindRandomFact(ctx context.Context) (Fact, error)
}

func (h randomFactHandler) Handle(ctx context.Context, query RandomFact) (f Fact, err error) {
	return h.readModel.FindRandomFact(ctx)
}
