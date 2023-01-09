package main

import (
	"os"

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

func setupDatabaseHandler(connectionString string, mongoDatabaseName string) *database.DatabaseHandler {
	databaseHandler := database.NewDatabaseHandler(connectionString, mongoDatabaseName)
	return &databaseHandler
}

func setupFactHandler(databaseHandler *database.DatabaseHandler) *facts.FactHandler {
	factHandler := facts.NewFactHandler(*databaseHandler)
	return &factHandler
}

func setupRouter(factHandler *facts.FactHandler) *router.Router {
	router := router.NewRouter(*factHandler)
	return &router
}

// @title           Animal Facts API
// @version         1.0
// @description     Get awesome facts about animals.
// @termsOfService  https://animalfact.app/terms

// @contact.name   Animal Facts API
// @contact.url    https://animalfacts.app/support
// @contact.email  support@animalfacts.app

// @license.name  MIT
// @license.url   https://github.com/cafo13/animal-facts/blob/main/LICENSE

// @host      https://animalfacts.app
// @BasePath  /api/v1
func main() {
	setupLogger()

	mongoDatabaseConnectionString := getEnvVar("MONGODB_CONNSTRING", "mongodb://animalfacts:animalfacts@localhost:27017")
	mongoDatabaseName := getEnvVar("MONGODB_DBNAME", "animalfacts")
	apiPort := getEnvVar("API_PORT", "8080")

	databaseHandler := setupDatabaseHandler(mongoDatabaseConnectionString, mongoDatabaseName)
	factHandler := setupFactHandler(databaseHandler)
	router := setupRouter(factHandler)

	router.StartRouter(apiPort)
}
