package main

import (
	"os"

	"github.com/cafo13/animal-facts/api/auth"
	"github.com/cafo13/animal-facts/api/database"
	"github.com/cafo13/animal-facts/api/facts"
	"github.com/cafo13/animal-facts/api/router"

	log "github.com/sirupsen/logrus"
)

func getEnvVar(envVar string, defaultValue string) string {
	if value, exists := os.LookupEnv(envVar); exists {
		return value
	}
	return defaultValue
}

func setupLogger() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func setupDatabaseHandler(dbHost string, dbPort string, dbName string, dbUser string, dbPassword string) (*database.DatabaseHandler, error) {
	databaseHandler, err := database.NewDatabaseHandler(dbHost, dbPort, dbName, dbUser, dbPassword)
	if err != nil {
		return nil, err
	}
	return &databaseHandler, nil
}

func setupAuthHandler(databaseHandler *database.DatabaseHandler) *auth.AuthHandler {
	authHandler := auth.NewAuthHandler(*databaseHandler)
	return &authHandler
}

func setupFactHandler(databaseHandler *database.DatabaseHandler) *facts.FactHandler {
	factHandler := facts.NewFactHandler(*databaseHandler)
	return &factHandler
}

func setupRouter(authHandler *auth.AuthHandler, factHandler *facts.FactHandler) router.GinRouter {
	router := router.NewRouter(*authHandler, *factHandler)
	return router
}

// @title           Animal Facts API
// @version         0.0.1
// @description     Awesome facts about animals.

// @contact.name   Animal Facts API
// @contact.url    https://animalfacts.app

// @license.name  MIT
// @license.url   https://github.com/cafo13/animal-facts/blob/main/LICENSE

// @host      https://animalfacts.app
// @BasePath  /api/v1
func main() {
	setupLogger()

	pgHost := getEnvVar("DB_HOST", "localhost")
	pgPort := getEnvVar("DB_PORT", "5432")
	pgDatabase := getEnvVar("DB_NAME", "animalfacts")
	pgUsername := getEnvVar("DB_USERNAME", "animalfacts")
	pgPassword := getEnvVar("DB_PASSWORD", "animalfacts")
	apiPort := getEnvVar("API_PORT", "8080")

	databaseHandler, err := setupDatabaseHandler(pgHost, pgPort, pgDatabase, pgUsername, pgPassword)
	if err != nil {
		log.Fatal("failed to setup database handler", err)
	}

	authHandler := setupAuthHandler(databaseHandler)
	factHandler := setupFactHandler(databaseHandler)
	router := setupRouter(authHandler, factHandler)

	router.StartRouter(apiPort)
}
