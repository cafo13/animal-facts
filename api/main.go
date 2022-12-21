package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cafo13/animal-facts/api/database"
	"github.com/cafo13/animal-facts/api/facts"
	"github.com/cafo13/animal-facts/api/types"

	"github.com/gin-gonic/gin"
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

func setupRouter(databaseHandler *database.DatabaseHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(http.StatusOK, "healthy\n")
	})

	router.GET("/fact", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")

		var fact *types.Fact

		factHandler, err := facts.NewFactHandler(*databaseHandler)
		if err != nil {
			log.Error("Error on initializing fact handler", err)
		}

		fact, err = factHandler.GetRandomFact()
		if err != nil {
			fmt.Println("Error on getting random fact", err)
			c.JSON(http.StatusInternalServerError, gin.H{"fact": types.Fact{}, "error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"fact": fact})
		}
	})

	router.GET("/fact/:id", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")

		var fact *types.Fact
		id := c.Params.ByName("id")

		factHandler, err := facts.NewFactHandler(*databaseHandler)
		if err != nil {
			log.Error("Error on initializing fact handler", err)
		}

		fact, err = factHandler.GetFactById(id)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error on getting fact by id %s", id), err)
			c.JSON(http.StatusNotFound, gin.H{"fact": types.Fact{}, "error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"fact": fact})
		}
	})

	return router
}

func main() {
	setupLogger()

	mongoDatabaseConnectionString := getEnvVar("MONGODB_CONNSTRING", "mongodb://animalfacts:animalfacts@localhost:27017")
	mongoDatabaseName := getEnvVar("MONGODB_DBNAME", "animalfacts")
	apiPort := getEnvVar("API_PORT", "8080")

	router := setupRouter(setupDatabaseHandler(mongoDatabaseConnectionString, mongoDatabaseName))
	router.Run(":" + apiPort)
}
