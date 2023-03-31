package main

import (
	"context"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"

	"github.com/cafo13/animal-facts/backend/facts-api/auth"
	"github.com/cafo13/animal-facts/backend/facts-api/repository"
	"github.com/cafo13/animal-facts/backend/facts-api/router"

	firebase "firebase.google.com/go/v4"
	log "github.com/sirupsen/logrus"
)

func setupLogger() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func setupAuthMiddleware() *auth.AuthMiddleware {
	if mockAuth, _ := strconv.ParseBool(os.Getenv("MOCK_AUTH")); mockAuth {
		authMiddleware := auth.NewMockAuthMiddleware()
		return &authMiddleware
	}

	config := &firebase.Config{ProjectID: os.Getenv("GCP_PROJECT")}
	firebaseApp, err := firebase.NewApp(context.Background(), config)
	if err != nil {
		panic(err)
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		panic(err)
	}

	authMiddleware := auth.NewFirebaseAuthMiddleware(authClient)
	return &authMiddleware
}

func setupFactsRepository(ctx context.Context, gcpProject string) *repository.FactsRepository {
	client, err := firestore.NewClient(ctx, gcpProject)
	if err != nil {
		panic(err)
	}

	factsFirestoreRepository := repository.NewFactsFirestoreRepository(client)
	return &factsFirestoreRepository
}

func setupRouter(authHandler *auth.AuthMiddleware, factsRepository *repository.FactsRepository) router.GinRouter {
	router := router.NewRouter(*authHandler, *factsRepository)
	return router
}

func main() {
	setupLogger()

	apiPort := os.Getenv("API_PORT")

	authMiddleware := setupAuthMiddleware()
	factsRepository := setupFactsRepository(context.Background(), os.Getenv("GCP_PROJECT"))
	router := setupRouter(authMiddleware, factsRepository)

	router.StartRouter(apiPort)
}