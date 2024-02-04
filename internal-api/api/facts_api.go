package api

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/cafo13/animal-facts/internal-api/docs"
	"github.com/cafo13/animal-facts/internal-api/handler"
	"github.com/cafo13/animal-facts/pkg/middleware"
	"github.com/cafo13/animal-facts/pkg/router"
)

var (
	basePathV1 = "api/v1"
)

type CreateFactResult struct {
	Id string `json:"id"`
}

type CreateUpdateFact struct {
	Fact   string `json:"fact"`
	Source string `json:"source"`
}

type ErrorResult struct {
	Error string `json:"error"`
}

type CountResult struct {
	Count int `json:"count"`
}

type FactsApi struct {
	factsApiRoutes []router.Route
	factsHandler   *handler.FactsHandler
}

func NewFactsApi(factsHandler *handler.FactsHandler) *FactsApi {
	return &FactsApi{factsHandler: factsHandler}
}

func (f *FactsApi) SetupRoutes() {
	f.factsApiRoutes = []router.Route{
		{
			Method:      "GET",
			Path:        "/swagger/*",
			HandlerFunc: echoSwagger.WrapHandler,
		},
		{
			Method:      "GET",
			Path:        "/health-internal",
			HandlerFunc: f.getHealth,
		},
		{
			Method:      "GET",
			Path:        fmt.Sprintf("/%s/facts/all", basePathV1),
			HandlerFunc: f.getAllFacts,
			Middlewares: []echo.MiddlewareFunc{
				middleware.EnsureValidToken(),
				middleware.VerifyScope("get:fact"),
			},
		},
		{
			Method:      "POST",
			Path:        fmt.Sprintf("/%s/facts", basePathV1),
			HandlerFunc: f.createFact,
			Middlewares: []echo.MiddlewareFunc{
				middleware.EnsureValidToken(),
				middleware.VerifyScope("create:fact"),
			},
		},
		{
			Method:      "PUT",
			Path:        fmt.Sprintf("%s/facts/:id", basePathV1),
			HandlerFunc: f.updateFact,
			Middlewares: []echo.MiddlewareFunc{
				middleware.EnsureValidToken(),
				middleware.VerifyScope("update:fact"),
			},
		},
		{
			Method:      "DELETE",
			Path:        fmt.Sprintf("%s/facts/:id", basePathV1),
			HandlerFunc: f.deleteFact,
			Middlewares: []echo.MiddlewareFunc{
				middleware.EnsureValidToken(),
				middleware.VerifyScope("delete:fact"),
			},
		},
		{
			Method:      "POST",
			Path:        fmt.Sprintf("/%s/facts/:id/approve", basePathV1),
			HandlerFunc: f.approveFact,
			Middlewares: []echo.MiddlewareFunc{
				middleware.EnsureValidToken(),
				middleware.VerifyScope("approve:fact"),
			},
		},
		{
			Method:      "POST",
			Path:        fmt.Sprintf("/%s/facts/:id/unapprove", basePathV1),
			HandlerFunc: f.unapproveFact,
			Middlewares: []echo.MiddlewareFunc{
				middleware.EnsureValidToken(),
				middleware.VerifyScope("unapprove:fact"),
			},
		},
	}
}

func (f *FactsApi) GetRoutes() []router.Route {
	return f.factsApiRoutes
}

func (f *FactsApi) getHealth(c echo.Context) error {
	// TODO check database connection for health check and maybe other connections?
	return c.String(http.StatusOK, "Healthy")
}

// createFact
//
//	@Summary      create fact
//	@Description  create a new fact
//	@Produce      json
//	@Param        request body CreateUpdateFact true "fact"
//	@Success      201  {object}  CreateFactResult
//	@Failure      400  {object}  ErrorResult
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts [post]
func (f *FactsApi) createFact(c echo.Context) error {
	fact := &CreateUpdateFact{}
	if err := c.Bind(fact); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResult{Error: err.Error()})
	}

	id := primitive.NewObjectID()
	err := f.factsHandler.Create(&handler.Fact{
		ID:       id,
		Fact:     fact.Fact,
		Source:   fact.Source,
		Approved: false,
	})
	if err != nil {
		// TODO only log error and return generic message as internal server error should not be displayed to user
		return c.JSON(http.StatusInternalServerError, ErrorResult{Error: err.Error()})
	}

	t := id.Hex()
	return c.JSON(http.StatusCreated, CreateFactResult{Id: t})
}

// updateFact
//
//	@Summary      update fact
//	@Description  update an existing fact
//	@Produce      json
//	@Param        request body CreateUpdateFact true "fact"
//	@Success      200  {string}  "fact updated"
//	@Failure      400  {object}  ErrorResult
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts/:id [put]
func (f *FactsApi) updateFact(c echo.Context) error {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResult{Error: "id from request path is not a valid object id in hex string format"})
	}

	fact := &CreateUpdateFact{}
	if err := c.Bind(fact); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResult{Error: err.Error()})
	}

	err = f.factsHandler.Update(&handler.Fact{
		ID:     objID,
		Fact:   fact.Fact,
		Source: fact.Source,
	})
	if err != nil {
		// TODO only log error and return generic message as internal server error should not be displayed to user
		return c.JSON(http.StatusInternalServerError, ErrorResult{Error: err.Error()})
	}

	return c.String(http.StatusOK, "fact updated")
}

// deleteFact
//
//	@Summary      delete fact
//	@Description  delete an existing fact
//	@Produce      json
//	@Success      200  {string}  "fact deleted"
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts/:id [delete]
func (f *FactsApi) deleteFact(c echo.Context) error {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResult{Error: "id from request path is not a valid object id in hex string format"})
	}

	err = f.factsHandler.Delete(objID)
	if err != nil {
		// TODO only log error and return generic message as internal server error should not be displayed to user
		return c.JSON(http.StatusInternalServerError, ErrorResult{Error: err.Error()})
	}

	return c.String(http.StatusOK, "fact deleted")
}

// approveFact
//
//	@Summary      approve fact
//	@Description  approve an existing fact, so that it gets available in the public API
//	@Produce      json
//	@Success      200  {string}  "fact approved"
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts/:id/approve [post]
func (f *FactsApi) approveFact(c echo.Context) error {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResult{Error: "id from request path is not a valid object id in hex string format"})
	}

	err = f.factsHandler.Approve(objID)
	if err != nil {
		// TODO only log error and return generic message as internal server error should not be displayed to user
		return c.JSON(http.StatusInternalServerError, ErrorResult{Error: err.Error()})
	}

	return c.String(http.StatusOK, "fact approved")
}

// unapproveFact
//
//	@Summary      unapprove fact
//	@Description  unapprove an existing fact, so that it is no longer available in the public API
//	@Produce      json
//	@Success      200  {string}  "fact unapproved"
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts/:id/unapprove [post]
func (f *FactsApi) unapproveFact(c echo.Context) error {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResult{Error: "id from request path is not a valid object id in hex string format"})
	}

	err = f.factsHandler.Unapprove(objID)
	if err != nil {
		// TODO only log error and return generic message as internal server error should not be displayed to user
		return c.JSON(http.StatusInternalServerError, ErrorResult{Error: err.Error()})
	}

	return c.String(http.StatusOK, "fact unapproved")
}

// getAllFacts
//
//	@Summary      gets all facts
//	@Description  gets all facts (approved and unapproved) from the database
//	@Produce      json
//	@Success      200  {array}   []handler.Fact
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts/all     [get]
func (f *FactsApi) getAllFacts(c echo.Context) error {
	facts, err := f.factsHandler.GetAll()
	if err != nil {
		// TODO only log error and return generic message as internal server error should not be displayed to user
		return c.JSON(http.StatusInternalServerError, ErrorResult{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, &facts)
}
