package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/cafo13/animal-facts/pkg/router"
	"github.com/cafo13/animal-facts/public-api/handler"
)

var (
	basePathV1 = "api/v1"
)

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
			Path:        fmt.Sprintf("%s/facts", basePathV1),
			HandlerFunc: f.GetRandomApproved,
		},
		{
			Method:      "GET",
			Path:        fmt.Sprintf("%s/facts/:id", basePathV1),
			HandlerFunc: f.Get,
		},
	}
}

func (f *FactsApi) GetRoutes() []router.Route {
	return f.factsApiRoutes
}

func (f *FactsApi) GetRandomApproved(c echo.Context) error {
	return errors.New("not implemented")
}

func (f *FactsApi) Get(c echo.Context) error {
	return nil
}
