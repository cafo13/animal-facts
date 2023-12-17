//go:build integration

package integration_test

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/cafo13/animal-facts/public-api/server"
)

func Test_RunIntegrationTests(t *testing.T) {
	// making sure that API_PORT and MONGODB_URI are set for integration test setup
	apiPort, ok := os.LookupEnv("API_PORT")
	if !ok {
		apiPort = "8080"
	}

	err := os.Setenv("MONGODB_DATABASE", "animal-facts-integration-test")
	if err != nil {
		t.Error(err, "failed to set MONGODB_DATABASE env var for integration test runs")
		return
	}

	err = startPublicAPIServer(apiPort)
	if err != nil {
		t.Error(err, "failed to start public api for integration test runs")
		return
	}

	tests := []struct {
		name           string
		requestPath    string
		wantResponse   string
		wantHttpStatus int
		wantErr        bool
	}{
		{
			name:           "get count of facts in test database",
			requestPath:    "api/v1/facts/count",
			wantResponse:   `{"count":3}`,
			wantHttpStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "get fact by id from test database",
			requestPath:    "api/v1/facts/6578bf140e487ecc049c7594",
			wantResponse:   `{"id":"6578bf140e487ecc049c7594","fact":"The Blue Whale is the largest animal that has ever lived.","source":"https://factanimal.com/blue-whale/"}`,
			wantHttpStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "get random fact from test database",
			requestPath:    "api/v1/facts",
			wantHttpStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "get not found status on trying to get fact with id that not exists in test database",
			requestPath:    "api/v1/facts/507f1f77bcf86cd799439011",
			wantResponse:   `{"error":"fact with ID '507f1f77bcf86cd799439011' not found"}`,
			wantHttpStatus: http.StatusNotFound,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%s/%s", apiPort, tt.requestPath), nil)
			resp, _ := http.DefaultClient.Do(req)

			if resp.StatusCode != tt.wantHttpStatus {
				t.Errorf("got http status %d, want %d", resp.StatusCode, tt.wantHttpStatus)
				return
			}

			if tt.wantResponse != "" {
				body := make([]byte, resp.ContentLength)
				_, _ = resp.Body.Read(body)
				stringBody := strings.TrimSuffix(string(body), "\n")
				if !reflect.DeepEqual(stringBody, tt.wantResponse) {
					t.Errorf("got response %s, want %s", string(body), tt.wantResponse)
					return
				}
			}

			_ = resp.Body.Close()
		})
	}
}

func startPublicAPIServer(apiPort string) error {
	go server.Run()

	timeout := time.After(60 * time.Second)
	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-timeout:
			return errors.New("timed out while waiting for public api server to return http status 200 at /Health")
		case <-tick:
			req, _ := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:%s/health", apiPort), nil)
			resp, _ := http.DefaultClient.Do(req)
			if resp != nil {
				if resp.StatusCode == http.StatusOK {
					return nil
				}
				_ = resp.Body.Close()
			}
		}
	}
}
