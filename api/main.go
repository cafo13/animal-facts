package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cafo13/animal-facts/api/database"
	"github.com/cafo13/animal-facts/api/facts"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
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

func setupDatabaseConnection(connectionString string, mongoDatabaseName string) (*mongo.Database, error) {
	database, err := database.GetDatabase(connectionString, mongoDatabaseName)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to get database by connection %s", connectionString), err)
		return nil, err
	}
	return database, nil
}

func setupRouter(database mongo.Database) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "healthy\n")
	})

	router.GET("/fact", func(c *gin.Context) {
		var fact facts.Fact

		factHandler, err := facts.NewFactHandler(facts.Database{Db: database})
		if err != nil {
			log.Error("Error on initializing fact handler", err)
		}

		fact, err = factHandler.GetRandomFact()
		if err != nil {
			fmt.Println("Error on getting random fact", err)
			c.JSON(http.StatusNotFound, gin.H{"fact": "random fact not found", "error": err})
		}

		c.JSON(http.StatusOK, gin.H{"fact": fact})
	})

	router.GET("/fact/:id", func(c *gin.Context) {
		var fact facts.Fact
		id := c.Params.ByName("id")

		factHandler, err := facts.NewFactHandler(facts.Database{Db: database})
		if err != nil {
			log.Error("Error on initializing fact handler", err)
		}

		fact, err = factHandler.GetFactById(id)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error on getting fact handler by id %s", id), err)
			c.JSON(http.StatusNotFound, gin.H{"fact": fmt.Sprintf("fact with id %s not found", id), "error": err})
		}

		c.JSON(http.StatusOK, gin.H{"fact": fact})
	})

	return router
}

func main() {
	setupLogger()

	mongoDatabaseConnectionString := getEnvVar("MONGODB_CONNSTRING", "mongodb://animalfacts:animalfacts@localhost:27017")
	mongoDatabaseName := getEnvVar("MONGODB_DBNAME", "animalfacts")
	apiPort := getEnvVar("API_PORT", "8080")

	database, err := setupDatabaseConnection(mongoDatabaseConnectionString, mongoDatabaseName)
	if err != nil {
		log.Fatal("Failed to setup database connection", err)
	}

	router := setupRouter(*database)
	router.Run(":" + apiPort)
}
