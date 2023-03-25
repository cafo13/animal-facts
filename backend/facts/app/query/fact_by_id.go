package query

import (
	"context"

	"github.com/cafo13/animal-facts/backend/common/decorator"
	"github.com/sirupsen/logrus"
)

type FactByID struct {
	UUID string
}

type FactByIDHandler decorator.QueryHandler[FactByID, Fact]

type factByIDHandler struct {
	readModel FactByIDReadModel
}

func NewFactByIDHandler(
	readModel FactByIDReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) FactByIDHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[FactByID, Fact](
		factByIDHandler{readModel: readModel},
		logger,
		metricsClient,
	)
}

type FactByIDReadModel interface {
	FindFactByID(ctx context.Context, factUUID string) (Fact, error)
}

func (h factByIDHandler) Handle(ctx context.Context, query FactByID) (f Fact, err error) {
	return h.readModel.FindFactByID(ctx, query.UUID)
}
