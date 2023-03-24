package main

import (
	"context"
	"net/http"

	"github.com/cafo13/animal-facts/backend/common/logs"
	"github.com/cafo13/animal-facts/backend/common/server"
	"github.com/cafo13/animal-facts/backend/facts/openapi"
	"github.com/cafo13/animal-facts/backend/facts/service"

	chi "github.com/go-chi/chi/v5"
)

func main() {
	logs.Init()

	ctx := context.Background()

	app, cleanup := service.NewApplication(ctx)
	defer cleanup()

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return openapi.HandlerFromMux(openapi.NewHttpServer(app), router)
	})
}
