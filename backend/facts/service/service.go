package service

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	grpcClient "github.com/cafo13/animal-facts/backend/common/client"
	"github.com/cafo13/animal-facts/backend/common/metrics"
	"github.com/cafo13/animal-facts/backend/trainings/adapters"
	"github.com/cafo13/animal-facts/backend/trainings/app"
	"github.com/cafo13/animal-facts/backend/trainings/app/command"
	"github.com/cafo13/animal-facts/backend/trainings/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	trainerClient, closeTrainerClient, err := grpcClient.NewTrainerClient()
	if err != nil {
		panic(err)
	}

	usersClient, closeUsersClient, err := grpcClient.NewUsersClient()
	if err != nil {
		panic(err)
	}
	trainerGrpc := adapters.NewTrainerGrpc(trainerClient)
	usersGrpc := adapters.NewUsersGrpc(usersClient)

	return newApplication(ctx, trainerGrpc, usersGrpc),
		func() {
			_ = closeTrainerClient()
			_ = closeUsersClient()
		}
}

func NewComponentTestApplication(ctx context.Context) app.Application {
	return newApplication(ctx, TrainerServiceMock{}, UserServiceMock{})
}

func newApplication(ctx context.Context, trainerGrpc command.TrainerService, usersGrpc command.UserService) app.Application {
	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	trainingsRepository := adapters.NewTrainingsFirestoreRepository(client)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return app.Application{
		Commands: app.Commands{
			ApproveTrainingReschedule: command.NewApproveTrainingRescheduleHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
			CancelTraining:            command.NewCancelTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
			RejectTrainingReschedule:  command.NewRejectTrainingRescheduleHandler(trainingsRepository, logger, metricsClient),
			RescheduleTraining:        command.NewRescheduleTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
			RequestTrainingReschedule: command.NewRequestTrainingRescheduleHandler(trainingsRepository, logger, metricsClient),
			ScheduleTraining:          command.NewScheduleTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
		},
		Queries: app.Queries{
			AllTrainings:     query.NewAllTrainingsHandler(trainingsRepository, logger, metricsClient),
			TrainingsForUser: query.NewTrainingsForUserHandler(trainingsRepository, logger, metricsClient),
		},
	}
}
