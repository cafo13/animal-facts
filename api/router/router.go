package router

import (
	"fmt"
	"net/http"

	"github.com/cafo13/animal-facts/api/docs"
	"github.com/cafo13/animal-facts/api/facts"
	"github.com/cafo13/animal-facts/api/types"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinRouter interface {
	GetHealth(context *gin.Context)
	GetRandomFact(context *gin.Context)
	GetFactById(context *gin.Context)
	StartRouter(port string)
}

type Router struct {
	Router      *gin.Engine
	FactHandler facts.FactHandler
}

func NewRouter(factHandler facts.FactHandler) Router {
	return Router{Router: gin.Default(), FactHandler: factHandler}
}

// @Summary Get health status
// @Schemes https
// @Description Checking the health of the API
// @Tags general
// @Produce plain
// @Success 200 {string} healthy
// @Router /healthy [get]
func (r Router) GetHealth(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "GET")
	context.String(http.StatusOK, "healthy\n")
}

// @Summary Get random animal fact
// @Schemes https
// @Description Getting a random animal fact
// @Tags facts
// @Produce json
// @Success 200 {object} types.Fact
// @Router /fact [get]
func (r Router) GetRandomFact(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "GET")

	var fact *types.Fact

	fact, err := r.FactHandler.GetRandomFact()
	if err != nil {
		fmt.Println("Error on getting random fact", err)
		context.JSON(http.StatusInternalServerError, gin.H{"Fact": types.Fact{}, "error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"Fact": fact})
	}
}

// @Summary Get animal fact by ID
// @Schemes https
// @Description Getting an animal fact by ID
// @Tags facts
// @Produce json
// @Success 200 {object} types.Fact
// @Router /fact/:id [get]
func (r Router) GetFactById(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "GET")

	var fact *types.Fact
	id := context.Params.ByName("id")

	fact, err := r.FactHandler.GetFactById(id)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error on getting fact by id %s", id), err)
		context.JSON(http.StatusNotFound, gin.H{"Fact": types.Fact{}, "error": err.Error()})
	} else {
		context.JSON(http.StatusOK, gin.H{"Fact": fact})
	}
}

func (r Router) StartRouter(port string) {
	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := r.Router.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			r.GetHealth(c)
		})

		v1.GET("/fact", func(c *gin.Context) {
			r.GetRandomFact(c)
		})

		v1.GET("/fact/:id", func(c *gin.Context) {
			r.GetFactById(c)
		})
	}

	r.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Router.Run(":" + port)
}
