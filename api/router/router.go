package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cafo13/animal-facts/api/database"
	"github.com/cafo13/animal-facts/api/docs"
	"github.com/cafo13/animal-facts/api/facts"

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
// @Success 200 {object} database.Fact
// @Failure 500 {object} ErrorResponse
// @Router /fact [get]
func (r Router) GetRandomFact(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "GET")

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
// @Success 200 {object} database.Fact
// @Failure 404 {object} ErrorResponse
// @Router /fact/:id [get]
func (r Router) GetFactById(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "GET")

	id := context.Params.ByName("id")
	uintId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		wrappedError := errors.Wrap(err, fmt.Sprintf("error on parsing id %s", id))
		log.Error(wrappedError)
		context.JSON(http.StatusNotFound, gin.H{"Error": wrappedError.Error()})
		return
	}

	fact, err := r.FactHandler.GetFactById(uint(uintId))
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
// @Param request body database.Fact true "a new fact"
// @Produce json
// @Success 201 {object} database.Fact
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /fact [post]
func (r Router) AddFact(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "POST")

	fact := &database.Fact{}
	err := context.BindJSON(&fact)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting fact from json body")
		log.Error(wrappedError)
		context.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError.Error()})
		return
	}

	err = r.FactHandler.CreateFact(fact)
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
// @Success 200 {object} database.Fact
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} MessageResponse
// @Failure 500 {object} ErrorResponse
// @Router /fact/:id [put]
func (r Router) UpdateFact(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "PUT")

	id := context.Params.ByName("id")
	uintId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		wrappedError := errors.Wrap(err, fmt.Sprintf("error on parsing id %s", id))
		log.Error(wrappedError)
		context.JSON(http.StatusNotFound, gin.H{"Error": wrappedError.Error()})
		return
	}

	fact := &database.Fact{}
	fact.ID = uint(uintId)
	err = fact.Read()
	if err != nil {
		errorMsg := fmt.Sprintf("error on updating fact, fact with id %s does not exists", id)
		log.Error(errorMsg)
		context.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	err = context.BindJSON(&fact)

	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting fact from json body")
		log.Error(wrappedError)
		context.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError})
		return
	}

	updatedFact, err := r.FactHandler.UpdateFact(uint(uintId), fact)
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
	uintId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		wrappedError := errors.Wrap(err, fmt.Sprintf("error on parsing id %s", id))
		log.Error(wrappedError)
		context.JSON(http.StatusNotFound, gin.H{"Error": wrappedError.Error()})
		return
	}

	err = r.FactHandler.DeleteFact(uint(uintId))
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
		v1.GET("/health", r.GetHealth)

		v1.GET("/fact", r.GetRandomFact)

		v1.GET("/fact/:id", r.GetFactById)

		v1.POST("/fact", r.AddFact)

		v1.PUT("/fact/:id", r.UpdateFact)

		v1.DELETE("/fact/:id", r.DeleteFact)
	}

	r.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Router.Run(":" + port)
}
