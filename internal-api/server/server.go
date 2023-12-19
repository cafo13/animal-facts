package server

import (
	"context"
	"os"
	"os/signal"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/pkg/errors"

	"github.com/cafo13/animal-facts/internal-api/api"
	"github.com/cafo13/animal-facts/internal-api/handler"
	logger "github.com/cafo13/animal-facts/pkg/log"
	"github.com/cafo13/animal-facts/pkg/repository"
	"github.com/cafo13/animal-facts/pkg/router"
	"github.com/cafo13/animal-facts/pkg/service"
)

var (
	mongoDbUri string
)

// Run
//
// @title           Animal Facts Internal API
// @version         VERSION_PLACEHOLDER
// @description     This API provides facts about animals.
//
// @license.name  MIT
// @license.url   https://github.com/cafo13/animal-facts/blob/main/LICENSE
//
// @host      https://animal-facts-manager.cafo.dev
// @BasePath  /api/v1
//
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func Run() {
	logger.SetupLogger()
	log.Logger().Info("starting internal animal facts api VERSION_PLACEHOLDER")

	loadEnv()

	factsRouter, err := setupServiceDependencies()
	if err != nil {
		panic(errors.Wrap(err, "failed to setup service dependencies"))
	}

	svc := service.NewService(factsRouter)

	apiPortStr, ok := os.LookupEnv("API_PORT")
	if !ok {
		apiPortStr = "8080"
		log.Logger().Infof("API_PORT environment variable is not set, using default value %s", apiPortStr)
	}

	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse API_PORT environment variable, only integer values are allowed (like 80 or 8080"))
	}

	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt)
	defer cancel()

	err = svc.Run(ctx, apiPort)
	if err != nil {
		panic(errors.Wrap(err, "failed to start service"))
	}
}

func loadEnv() {
	log.Logger().Info("loading environment variables")
	err := godotenv.Load()
	if err != nil {
		log.Logger().WithError(err).Warn("failed to load .env file")
	}

	var ok bool
	mongoDbUri, ok = os.LookupEnv("MONGODB_URI")
	if !ok {
		panic("MONGODB_URI environment variable is not set")
	}
}

func setupServiceDependencies() (*router.Router, error) {
	factsRepository, err := repository.NewMongoDBFactsRepository(mongoDbUri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup mongo db facts repository")
	}

	factsHandler := handler.NewFactsHandler(factsRepository)
	factsApi := api.NewFactsApi(factsHandler)
	factsApi.SetupRoutes()
	factsRouter := router.NewRouter()
	for _, route := range factsApi.GetRoutes() {
		err := factsRouter.RegisterRoute(route)
		if err != nil {
			log.Logger().WithError(err).Errorf("failed to register route %s", route.Path)
		} else {
			log.Logger().Infof("registered route %s", route.Path)
		}
	}

	return factsRouter, nil
}