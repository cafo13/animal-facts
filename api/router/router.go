package router

import (
	"fmt"
	"net/http"

	"github.com/cafo13/animal-facts/api/docs"
	"github.com/cafo13/animal-facts/api/facts"
	"github.com/cafo13/animal-facts/api/types"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinRouter interface {
	AddFact(context *gin.Context)
	DeleteFact(context *gin.Context)
	GetHealth(context *gin.Context)
	GetRandomFact(context *gin.Context)
	GetFactById(context *gin.Context)
	UpdateFact(context *gin.Context)
	StartRouter(port string)
}

type Router struct {
	Router      *gin.Engine
	FactHandler facts.FactHandler
}

type MessageResponse struct {
	Message string
}

type ErrorResponse struct {
	Error error
}

func NewRouter(factHandler facts.FactHandler) GinRouter {
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
// @Failure 500 {object} ErrorResponse
// @Router /fact [get]
func (r Router) GetRandomFact(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "GET")

	var fact *types.Fact

	fact, err := r.FactHandler.GetRandomFact()
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting random fact")
		log.Error(wrappedError)
		context.JSON(http.StatusInternalServerError, gin.H{"Error": wrappedError.Error()})
		return
	} else {
		context.IndentedJSON(http.StatusOK, fact)
		return
	}
}

// @Summary Get animal fact by id
// @Schemes https
// @Description Getting an animal fact by id
// @Tags facts
// @Produce json
// @Success 200 {object} types.Fact
// @Failure 404 {object} ErrorResponse
// @Router /fact/:id [get]
func (r Router) GetFactById(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "GET")

	var fact *types.Fact
	id := context.Params.ByName("id")

	fact, err := r.FactHandler.GetFactById(id)
	if err != nil {
		wrappedError := errors.Wrap(err, fmt.Sprintf("error on getting fact by id %s", id))
		log.Error(wrappedError)
		context.JSON(http.StatusNotFound, gin.H{"Error": wrappedError.Error()})
		return
	} else {
		context.IndentedJSON(http.StatusOK, fact)
		return
	}
}

// @Summary Add a new animal fact
// @Schemes https
// @Description Adding an animal fact
// @Tags facts
// @Accept json
// @Param request body types.Fact true "a new fact"
// @Produce json
// @Success 201 {object} types.Fact
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /fact [post]
func (r Router) AddFact(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "POST")

	var fact *types.Fact
	err := context.BindJSON(&fact)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting fact from json body")
		log.Error(wrappedError)
		context.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError.Error()})
		return
	}

	err = r.FactHandler.AddFact(fact)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on adding new fact")
		log.Error(wrappedError)
		context.JSON(http.StatusInternalServerError, gin.H{"Error": wrappedError.Error()})
		return
	} else {
		context.IndentedJSON(http.StatusCreated, fact)
		return
	}
}

// @Summary Update an existing animal fact
// @Schemes https
// @Description Updating an animal fact
// @Tags facts
// @Accept json
// @Param request body interface{} true "an updated fact"
// @Produce json
// @Success 200 {object} types.Fact
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} MessageResponse
// @Failure 500 {object} ErrorResponse
// @Router /fact/:id [put]
func (r Router) UpdateFact(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "PUT")

	id := context.Params.ByName("id")
	exists, err := r.FactHandler.FactExists(id)
	if err != nil {
		wrappedError := errors.Wrap(err, fmt.Sprintf("error on updating fact, error on checking, if fact with id %s exists", id))
		log.Error(wrappedError)
		context.JSON(http.StatusInternalServerError, gin.H{"Message": wrappedError.Error()})
		return
	} else if !exists {
		errorMsg := fmt.Sprintf("error on updating fact, fact with id %s does not exists", id)
		log.Error(errorMsg)
		context.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	var fact *interface{}
	err = context.BindJSON(&fact)

	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting fact from json body")
		log.Error(wrappedError)
		context.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError})
		return
	}

	updatedFact, err := r.FactHandler.UpdateFact(id, fact)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on updating fact")
		log.Error(wrappedError)
		context.JSON(http.StatusInternalServerError, gin.H{"Error": wrappedError})
		return
	} else {
		context.IndentedJSON(http.StatusOK, updatedFact)
		return
	}
}

// @Summary Delete an animal fact
// @Schemes https
// @Description Deleting an animal fact
// @Tags facts
// @Produce json
// @Success 200 {object} MessageResponse
// @Failure 500 {object} ErrorResponse
// @Router /fact/:id [delete]
func (r Router) DeleteFact(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "DELETE")

	id := context.Params.ByName("id")

	err := r.FactHandler.DeleteFact(id)
	if err != nil {
		wrappedError := errors.Wrap(err, fmt.Sprintf("error on deleting fact with id %s", id))
		log.Error(wrappedError)
		context.JSON(http.StatusInternalServerError, gin.H{"Error": wrappedError})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"Message": "deleted fact successfully"})
		return
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

		v1.POST("/fact", func(c *gin.Context) {
			r.AddFact(c)
		})

		v1.PUT("/fact/:id", func(c *gin.Context) {
			r.UpdateFact(c)
		})

		v1.DELETE("/fact/:id", func(c *gin.Context) {
			r.DeleteFact(c)
		})
	}

	r.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Router.Run(":" + port)
}
