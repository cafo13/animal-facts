package main

import (
	"github.com/joho/godotenv"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"

	"github.com/cafo13/animal-facts/internal/repository"
	"github.com/cafo13/animal-facts/internal/service"
)

var (
	mongoDbUri string
)

func main() {
	e := echo.New()

	setupLogger(e)

	loadEnv()

	factsRepository, err := repository.NewMongoDBFactsRepository(mongoDbUri)
	if err != nil {
		log.Logger().WithError(err).Fatal("failed to setup mongo db facts repository")
		os.Exit(1)
	}

	factsHandler, err := handler.NewFactsHandler(factsRepository)
	if err != nil {
		log.Logger().WithError(err).Fatal("failed to setup facts handler")
		os.Exit(1)
	}

	factsRouter, err := router.NewFactsRouter(factsHandler)
	if err != nil {
		log.Logger().WithError(err).Fatal("failed to setup facts router")
		os.Exit(1)
	}

	svc := service.NewService(factsRouter)
	err = svc.Start()
	if err != nil {
		log.Logger().WithError(err).Fatal("failed to start service")
		os.Exit(1)
	}
}

func setupLogger(e *echo.Echo) {
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	e.Logger = log.Logger()
	e.Use(middleware.Logger())
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
