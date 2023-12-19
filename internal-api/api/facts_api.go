package api

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"

	_ "github.com/cafo13/animal-facts/internal-api/docs"
	"github.com/cafo13/animal-facts/internal-api/handler"
	"github.com/cafo13/animal-facts/pkg/router"
)

var (
	basePathV1 = "api/v1"
)

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
			HandlerFunc: f.GetHealth,
		},
		{
			Method:      "POST",
			Path:        fmt.Sprintf("/%s/facts", basePathV1),
			HandlerFunc: f.CreateFact,
		},
		{
			Method:      "PUT",
			Path:        fmt.Sprintf("%s/facts/:id", basePathV1),
			HandlerFunc: f.UpdateFact,
		},
		{
			Method:      "DELETE",
			Path:        fmt.Sprintf("%s/facts/:id", basePathV1),
			HandlerFunc: f.DeleteFact,
		},
	}
}

func (f *FactsApi) GetRoutes() []router.Route {
	return f.factsApiRoutes
}

func (f *FactsApi) GetHealth(c echo.Context) error {
	// TODO check database connection for health check and maybe other connections?
	return c.String(http.StatusOK, "Healthy")
}

// CreateFact
//
//	@Summary      create fact
//	@Description  create a new fact
//	@Produce      json
//	@Param        request body CreateUpdateFact true "fact"
//	@Success      201  {string}  "fact created"
//	@Failure      400  {object}  ErrorResult
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts [post]
func (f *FactsApi) CreateFact(c echo.Context) error {
	fact := &CreateUpdateFact{}
	if err := c.Bind(fact); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResult{Error: err.Error()})
	}

	err := f.factsHandler.Create(&handler.Fact{
		ID:       primitive.NewObjectID(),
		Fact:     fact.Fact,
		Source:   fact.Source,
		Approved: fact.Approved,
	})
	if err != nil {
		// TODO only log error and return generic message as internal server error should not be displayed to user
		return c.JSON(http.StatusInternalServerError, ErrorResult{Error: err.Error()})
	}

	return c.String(http.StatusCreated, "fact created")
}

// UpdateFact
//
//	@Summary      update fact
//	@Description  update an existing fact
//	@Produce      json
//	@Param        request body CreateUpdateFact true "fact"
//	@Success      200  {string}  "fact updated"
//	@Failure      400  {object}  ErrorResult
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts/:id [put]
func (f *FactsApi) UpdateFact(c echo.Context) error {
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

// DeleteFact
//
//	@Summary      delete fact
//	@Description  delete an existing fact
//	@Produce      json
//	@Success      200  {string}  "fact deleted"
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts/:id [delete]
func (f *FactsApi) DeleteFact(c echo.Context) error {
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
