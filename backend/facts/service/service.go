package service

import (
	"context"
	"os"

	"github.com/cafo13/animal-facts/backend/common/metrics"
	"github.com/cafo13/animal-facts/backend/facts/adapters"
	"github.com/cafo13/animal-facts/backend/facts/app"
	"github.com/cafo13/animal-facts/backend/facts/app/command"
	"github.com/cafo13/animal-facts/backend/facts/app/query"

	"cloud.google.com/go/firestore"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	return newApplication(ctx),
		func() {
			// TODO: Execute necessary cleanups here
		}
}

func newApplication(ctx context.Context) app.Application {
	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	factsRepository := adapters.NewFactsFirestoreRepository(client)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return app.Application{
		Commands: app.Commands{
			CreateFact: command.NewCreateFactHandler(factsRepository, logger, metricsClient),
			UpdateFact: command.NewUpdateFactHandler(factsRepository, logger, metricsClient),
			DeleteFact: command.NewDeleteFactHandler(factsRepository, logger, metricsClient),
		},
		Queries: app.Queries{
			RandomFact: query.NewRandomFactHandler(factsRepository, logger, metricsClient),
		},
	}
}
