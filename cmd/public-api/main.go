package main

import (
	"context"
	"github.com/cafo13/animal-facts/public-api/api"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/cafo13/animal-facts/pkg/repository"
	"github.com/cafo13/animal-facts/pkg/router"
	"github.com/cafo13/animal-facts/public-api/handler"
	"github.com/cafo13/animal-facts/public-api/service"
)

var (
	mongoDbUri string
)

// @title           Animal Facts Public API
// @version         VERSION_PLACEHOLDER
// @description     This API provides facts about animals.

// @license.name  MIT
// @license.url   https://github.com/cafo13/animal-facts/blob/main/LICENSE

// @host      https://animal-facts.cafo.dev
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	setupLogger()
	log.Logger().Info("starting public animal facts api VERSION_PLACEHOLDER")

	loadEnv()

	factsRepository, err := repository.NewMongoDBFactsRepository(mongoDbUri)
	if err != nil {
		panic(errors.Wrap(err, "failed to setup mongo db facts repository"))
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

	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt)
	defer cancel()

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

	err = svc.Run(ctx, apiPort)
	if err != nil {
		panic(errors.Wrap(err, "failed to start service"))
	}
}

func setupLogger() {
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	log.Logger().Info("logger enabled")
}

func loadEnv() {
	log.Logger().Info("loading environment variables")
	err := godotenv.Load()
	if err != nil {
		log.Warn("failed to load .env file")
	}

	var ok bool
	mongoDbUri, ok = os.LookupEnv("MONGODB_URI")
	if !ok {
		log.Logger().Fatal("MONGODB_URI environment variable is not set")
		os.Exit(1)
	}
}
