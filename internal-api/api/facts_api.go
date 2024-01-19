package api

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/cafo13/animal-facts/internal-api/docs"
	"github.com/cafo13/animal-facts/internal-api/handler"
	"github.com/cafo13/animal-facts/pkg/router"
)

var (
	basePathV1 = "api/v1"
)

type CreateFactResult struct {
	Id string `json:"id"`
}

type CreateUpdateFact struct {
	Fact     string `json:"fact"`
	Source   string `json:"source"`
	Approved bool   `json:"approved"`
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
			Path:        "/health",
			HandlerFunc: f.getHealth,
		},
		{
			Method:      "POST",
			Path:        fmt.Sprintf("/%s/facts", basePathV1),
			HandlerFunc: f.createFact,
		},
		{
			Method:      "PUT",
			Path:        fmt.Sprintf("%s/facts/:id", basePathV1),
			HandlerFunc: f.updateFact,
		},
		{
			Method:      "DELETE",
			Path:        fmt.Sprintf("%s/facts/:id", basePathV1),
			HandlerFunc: f.deleteFact,
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
		Approved: fact.Approved,
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
		ID:       objID,
		Fact:     fact.Fact,
		Source:   fact.Source,
		Approved: fact.Approved,
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
