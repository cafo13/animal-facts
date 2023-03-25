package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cafo13/animal-facts/backend/auth"
	"github.com/cafo13/animal-facts/backend/repository"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type GinRouter interface {
	AddFact(ctx *gin.Context)
	DeleteFact(ctx *gin.Context)
	GetRandomFact(ctx *gin.Context)
	GetFactById(ctx *gin.Context)
	UpdateFact(ctx *gin.Context)

	StartRouter(port string)
}

type Router struct {
	Router          *gin.Engine
	AuthMiddleware  auth.AuthMiddleware
	FactsRepository repository.FactsRepository
}

type MessageResponse struct {
	Message string
}

type ErrorResponse struct {
	Error error
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewRouter(authMiddleware auth.AuthMiddleware, factsRepository repository.FactsRepository) GinRouter {
	return Router{Router: gin.Default(), AuthMiddleware: authMiddleware, FactsRepository: factsRepository}
}

func (r Router) AddFact(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "POST")

	fact := &repository.FactModel{}
	err := ctx.BindJSON(&fact)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting fact from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError.Error()})
		return
	}

	err = r.FactsRepository.CreateFact(ctx, fact)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusCreated, fact)
		return
	}
}

func (r Router) GetFactById(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET")

	uuid := ctx.Params.ByName("id")
	fact, err := r.FactsRepository.ReadFact(ctx, uuid)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, fact)
		return
	}
}

func (r Router) GetRandomFact(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET")

	fact, err := r.FactsRepository.ReadRandomFact(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, fact)
		return
	}
}

func (r Router) UpdateFact(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "PUT")

	uuid := ctx.Params.ByName("id")
	_, err := r.FactsRepository.ReadFact(ctx, uuid)
	if err != nil {
		errorMsg := fmt.Sprintf("error on loading fact with uuid %s", uuid)
		log.Error(errorMsg)
		ctx.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	var fact *repository.FactModel
	err = ctx.BindJSON(&fact)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting fact from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError})
		return
	}

	err = r.FactsRepository.UpdateFact(
		ctx,
		uuid,
		func(context context.Context, f *repository.FactModel) (*repository.FactModel, error) {
			if fact.Text != "" {
				f.Text = fact.Text
			}
			if fact.Source != "" {
				f.Source = fact.Source
			}
			if fact.Approved {
				f.Approved = fact.Approved
			}

			return f, nil
		},
	)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on updating fact")
		log.Error(wrappedError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": wrappedError})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"Message": "updated fact successfully"})
		return
	}
}

func (r Router) DeleteFact(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "DELETE")

	id := ctx.Params.ByName("id")
	err := r.FactsRepository.DeleteFact(ctx, id)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"Message": "deleted fact successfully"})
		return
	}
}

func (r Router) StartRouter(port string) {
	v1 := r.Router.Group("/api/v1")
	{
		v1.GET("/fact", r.GetRandomFact)

		v1.GET("/fact/:id", r.GetFactById)

		v1.POST("/fact", r.AuthMiddleware.Middleware(), r.AddFact)

		v1.PUT("/fact/:id", r.AuthMiddleware.Middleware(), r.UpdateFact)

		v1.DELETE("/fact/:id", r.AuthMiddleware.Middleware(), r.DeleteFact)
	}

	r.Router.Run(":" + port)
}
