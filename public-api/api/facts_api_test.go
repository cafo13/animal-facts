//go:build integration

package api

import (
	"github.com/pkg/errors"
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

func getFactsApi() (*FactsApi, error) {
	mongoDbUri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		return nil, errors.New("MONGODB_URI environment variable is not set, set it to a test database before running the integration tests")
	}

	fatsRepository, err := repository.NewMongoDBFactsRepository(mongoDbUri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup repository for integration tests")
	}

	factsHandler := handler.NewFactsHandler(fatsRepository)
	factsApi := NewFactsApi(factsHandler)
	return factsApi, nil
}

func Test_RunApiIntegrationTests_getCount(t *testing.T) {
	factsApi, err := getFactsApi()
	if err != nil {
		t.Error(err)
		return
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
	factsApi, err := getFactsApi()
	if err != nil {
		t.Error(err)
		return
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
	factsApi, err := getFactsApi()
	if err != nil {
		t.Error(err)
		return
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
