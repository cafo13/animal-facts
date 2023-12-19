//go:build integration

package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/cafo13/animal-facts/pkg/repository"
	"github.com/cafo13/animal-facts/public-api/handler"
	"github.com/labstack/echo/v4"
)

func Test_RunApiIntegrationTests_getCount(t *testing.T) {
	mongoDbUri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		t.Error("MONGODB_URI environment variable is not set, set it to a test database before running the integration tests")
	}

	tests := []struct {
		name           string
		wantResponse   string
		wantHttpStatus int
		wantErr        bool
	}{
		{
			name:           "get count of facts in test database",
			wantResponse:   `{"count":3}`,
			wantHttpStatus: http.StatusOK,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			fatsRepository, _ := repository.NewMongoDBFactsRepository(mongoDbUri)
			factsHandler := handler.NewFactsHandler(fatsRepository)
			factsApi := NewFactsApi(factsHandler)

			gotErr := factsApi.getCount(c)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("RunApiIntegrationTests_getCount() error = %v, wantErr = %v", gotErr, tt.wantErr)
				return
			}

			if rec.Code != tt.wantHttpStatus {
				t.Errorf("RunApiIntegrationTests_getCount() gotHttpStatus = %v, wantHttpStatus = %v", rec.Code, tt.wantHttpStatus)
			}

			if tt.wantResponse != "" {
				stringBody := strings.TrimSuffix(rec.Body.String(), "\n")
				if !reflect.DeepEqual(stringBody, tt.wantResponse) {
					t.Errorf("RunApiIntegrationTests_getCount() gotResponse = %v, wantResponse = %v", stringBody, tt.wantResponse)
				}
			}
		})
	}
}

func Test_RunApiIntegrationTests_get(t *testing.T) {
	mongoDbUri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		t.Error("MONGODB_URI environment variable is not set, set it to a test database before running the integration tests")
	}

	tests := []struct {
		name           string
		requestFactID  string
		wantResponse   string
		wantHttpStatus int
		wantErr        bool
	}{
		{
			name:           "get fact by id from test database",
			requestFactID:  "6578bf140e487ecc049c7594",
			wantResponse:   `{"id":"6578bf140e487ecc049c7594","fact":"The Blue Whale is the largest animal that has ever lived.","source":"https://factanimal.com/blue-whale/"}`,
			wantHttpStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "get not found status on trying to get fact with id that not exists in test database",
			requestFactID:  "507f1f77bcf86cd799439011",
			wantResponse:   `{"error":"fact with ID '507f1f77bcf86cd799439011' not found"}`,
			wantHttpStatus: http.StatusNotFound,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.requestFactID)

			fatsRepository, _ := repository.NewMongoDBFactsRepository(mongoDbUri)
			factsHandler := handler.NewFactsHandler(fatsRepository)
			factsApi := NewFactsApi(factsHandler)

			gotErr := factsApi.get(c)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("RunApiIntegrationTests_get() error = %v, wantErr = %v", gotErr, tt.wantErr)
				return
			}

			if rec.Code != tt.wantHttpStatus {
				t.Errorf("RunApiIntegrationTests_get() gotHttpStatus = %v, wantHttpStatus = %v", rec.Code, tt.wantHttpStatus)
			}

			if tt.wantResponse != "" {
				stringBody := strings.TrimSuffix(rec.Body.String(), "\n")
				if !reflect.DeepEqual(stringBody, tt.wantResponse) {
					t.Errorf("RunApiIntegrationTests_get() gotResponse = %v, wantResponse = %v", stringBody, tt.wantResponse)
				}
			}
		})
	}
}

func Test_RunApiIntegrationTests_getRandomApproved(t *testing.T) {
	mongoDbUri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		t.Error("MONGODB_URI environment variable is not set, set it to a test database before running the integration tests")
	}

	tests := []struct {
		name           string
		wantHttpStatus int
		wantErr        bool
	}{
		{
			name:           "get random fact from test database",
			wantHttpStatus: http.StatusOK,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			fatsRepository, _ := repository.NewMongoDBFactsRepository(mongoDbUri)
			factsHandler := handler.NewFactsHandler(fatsRepository)
			factsApi := NewFactsApi(factsHandler)

			gotErr := factsApi.getRandomApproved(c)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("RunApiIntegrationTests_getRandomApproved() error = %v, wantErr = %v", gotErr, tt.wantErr)
				return
			}

			if rec.Code != tt.wantHttpStatus {
				t.Errorf("RunApiIntegrationTests_getRandomApproved() gotHttpStatus = %v, wantHttpStatus = %v", rec.Code, tt.wantHttpStatus)
			}
		})
	}
}
