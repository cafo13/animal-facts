package main

import (
	"fmt"
	"net/http"
	"os"

	"animalfacts.app/handler"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setupLogger() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "healthy\n")
	})

	router.GET("/fact", func(c *gin.Context) {
		var fact handler.Fact

		factHandler, err := handler.NewFactHandler(handler.Database{Db: "test"})
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
		var fact handler.Fact
		id := c.Params.ByName("id")

		factHandler, err := handler.NewFactHandler(handler.Database{Db: "test"})
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
	router := setupRouter()
	router.Run(":8080")
}
