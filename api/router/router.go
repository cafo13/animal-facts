package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cafo13/animal-facts/api/auth"
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
	AuthHandler auth.AuthHandler
	FactHandler facts.FactHandler
}

type MessageResponse struct {
	Message string
}

type ErrorResponse struct {
	Error error
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewRouter(authHandler auth.AuthHandler, factHandler facts.FactHandler) GinRouter {
	return Router{Router: gin.Default(), AuthHandler: authHandler, FactHandler: factHandler}
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

// @Summary Login endpoint
// @Schemes https
// @Description Login endpoint for the endpoints of the API, that require auth
// @Tags general
// @Produce json
// @Success 200 {object}
// @Router /login [post]
func (r Router) Login(context *gin.Context) {
	var input LoginInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errors.Wrap(err, "invalid request format, needs username and password jso keys").Error()})
		return
	}

	u := database.User{}

	u.Username = input.Username
	u.Password = input.Password

	token, err := r.AuthHandler.VerifyLogin(u.Username, u.Password)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Get current user
// @Schemes https
// @Description Get current logged in user
// @Tags general
// @Produce json
// @Success 200 {object}
// @Router /user [get]
func (r Router) CurrentUser(context *gin.Context) {
	userId, err := auth.ExtractTokenID(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := r.AuthHandler.GetUserById(userId)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
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

	_, err = r.FactHandler.GetFactById(uint(uintId))
	if err != nil {
		errorMsg := fmt.Sprintf("error on updating fact, fact with id %s does not exists", id)
		log.Error(errorMsg)
		context.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	var fact *database.Fact
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

		v1.POST("/login", r.Login)

		v1.GET("/user", auth.JwtAuthMiddleware(), r.CurrentUser)

		v1.GET("/fact", r.GetRandomFact)

		v1.GET("/fact/:id", r.GetFactById)

		v1.POST("/fact", auth.JwtAuthMiddleware(), r.AddFact)

		v1.PUT("/fact/:id", auth.JwtAuthMiddleware(), r.UpdateFact)

		v1.DELETE("/fact/:id", auth.JwtAuthMiddleware(), r.DeleteFact)
	}

	r.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Router.Run(":" + port)
}
