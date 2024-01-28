package api

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/cafo13/animal-facts/pkg/router"
	_ "github.com/cafo13/animal-facts/public-api/docs"
	"github.com/cafo13/animal-facts/public-api/handler"
)

var (
	basePathV1 = "api/v1"
)

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
			Path:        "/swagger-public/*",
			HandlerFunc: echoSwagger.WrapHandler,
		},
		{
			Method:      "GET",
			Path:        "/health-public",
			HandlerFunc: f.getHealth,
		},
		{
			Method:      "GET",
			Path:        fmt.Sprintf("/%s/facts", basePathV1),
			HandlerFunc: f.getRandomApproved,
		},
		{
			Method:      "GET",
			Path:        fmt.Sprintf("%s/facts/:id", basePathV1),
			HandlerFunc: f.get,
		},
		{
			Method:      "GET",
			Path:        fmt.Sprintf("%s/facts/count", basePathV1),
			HandlerFunc: f.getCount,
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

// getRandomApproved
//
//	@Summary      gets random fact
//	@Description  gets random fact from the database
//	@Produce      json
//	@Success      200  {object}  handler.Fact
//	@Failure      404  {object}  ErrorResult
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts [get]
func (f *FactsApi) getRandomApproved(c echo.Context) error {
	fact, err := f.factsHandler.GetRandomApproved()
	if err != nil {
		// TODO only log error and return generic message as internal server error should not be displayed to user
		return c.JSON(http.StatusInternalServerError, ErrorResult{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, &fact)
}

// get
//
//	@Summary      gets fact
//	@Description  gets fact by ID from the database
//	@Produce      json
//	@Success      200  {object}  handler.Fact
//	@Failure      404  {object}  ErrorResult
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts/:id [get]
func (f *FactsApi) get(c echo.Context) error {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResult{Error: "id from request path is not a valid object id in hex string format"})
	}
	fact, err := f.factsHandler.Get(objID)
	if errors.Is(err, handler.ErrNotFound) {
		return c.JSON(http.StatusNotFound, ErrorResult{Error: fmt.Sprintf("fact with ID '%s' not found", id)})
	} else if err != nil {
		// TODO only log error and return generic message as internal server error should not be displayed to user
		return c.JSON(http.StatusInternalServerError, ErrorResult{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, &fact)
}

// getCount
//
//	@Summary      gets fact count
//	@Description  gets fact count from the database
//	@Produce      json
//	@Success      200  {object}  CountResult
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts/count [get]
func (f *FactsApi) getCount(c echo.Context) error {
	count, err := f.factsHandler.GetFactsCount()
	if err != nil {
		// TODO only log error and return generic message as internal server error should not be displayed to user
		return c.JSON(http.StatusInternalServerError, ErrorResult{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, &CountResult{Count: count})
}
