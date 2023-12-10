package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/swaggo/echo-swagger"

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
			Path:        "/swagger/*",
			HandlerFunc: echoSwagger.WrapHandler,
		},
		{
			Method:      "GET",
			Path:        "/health",
			HandlerFunc: f.GetHealth,
		},
		{
			Method:      "GET",
			Path:        fmt.Sprintf("/%s/facts", basePathV1),
			HandlerFunc: f.GetRandomApproved,
		},
		{
			Method:      "GET",
			Path:        fmt.Sprintf("%s/facts/:id", basePathV1),
			HandlerFunc: f.Get,
		},
		{
			Method:      "GET",
			Path:        fmt.Sprintf("%s/facts/count", basePathV1),
			HandlerFunc: f.GetCount,
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

// GetRandomApproved
//
//	@Summary      gets random fact
//	@Description  gets random fact from the database
//	@Produce      json
//	@Success      200  {object}  handler.Fact
//	@Failure      404  {object}  ErrorResult
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts [get]
func (f *FactsApi) GetRandomApproved(c echo.Context) error {
	return errors.New("not implemented")
}

// Get
//
//	@Summary      gets fact
//	@Description  gets fact by ID from the database
//	@Produce      json
//	@Success      200  {object}  handler.Fact
//	@Failure      404  {object}  ErrorResult
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts/:id [get]
func (f *FactsApi) Get(c echo.Context) error {
	return errors.New("not implemented")
}

// GetCount
//
//	@Summary      gets fact count
//	@Description  gets fact count from the database
//	@Produce      json
//	@Success      200  {object}  CountResult
//	@Failure      500  {object}  ErrorResult
//	@Router       /facts/count [get]
func (f *FactsApi) GetCount(c echo.Context) error {
	return errors.New("not implemented")
}
