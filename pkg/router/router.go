package router

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/neko-neko/echo-logrus/v2"
	"github.com/neko-neko/echo-logrus/v2/log"
	"net/http"
)

type Route struct {
	Method      string
	Path        string
	HandlerFunc func(c echo.Context) error
}

type Router struct {
	echoRouter *echo.Echo
}

func NewRouter() *Router {
	echoRouter := echo.New()
	echoRouter.Logger = log.Logger()
	echoRouter.Use(middleware.Logger())
	return &Router{echoRouter}
}

func (r Router) RegisterRoute(route Route) error {
	switch route.Method {
	case http.MethodGet:
		r.echoRouter.GET(route.Path, route.HandlerFunc)
		return nil
	case http.MethodPost:
		r.echoRouter.POST(route.Path, route.HandlerFunc)
		return nil
	case http.MethodPut:
		r.echoRouter.PUT(route.Path, route.HandlerFunc)
		return nil
	case http.MethodDelete:
		r.echoRouter.DELETE(route.Path, route.HandlerFunc)
		return nil
	}

	return fmt.Errorf("invalid http Method %s", route.Method)
}

func (r Router) Run(port int) error {
	return r.echoRouter.Start(fmt.Sprintf(":%v", port))
}

func (r Router) Shutdown(ctx context.Context) error {
	return r.echoRouter.Shutdown(ctx)
}