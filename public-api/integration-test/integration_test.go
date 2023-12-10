//go:build integration

package integration_test

import (
	"fmt"
	"net/http"
	"os"
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

	_, ok = os.LookupEnv("MONGODB_URI")
	if !ok {
		t.Error("MONGODB_URI env var needs to be set for integration tests, set it to a test database before running the tests")
		return
	}

	err := startPublicAPIServer(apiPort)
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
			requestPath:    "api/vs/facts/count",
			wantResponse:   `{ "count": 5 }`,
			wantHttpStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "get fact by id from test database",
			requestPath:    "api/vs/facts/3",
			wantResponse:   `{ "fact": "The Blue Whale is the largest animal that has ever lived.", "source": "https://factanimal.com/blue-whale/" }`,
			wantHttpStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "get random fact from test database",
			requestPath:    "api/vs/facts",
			wantHttpStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "get not found status on trying to get fact with id that not exists in test database",
			requestPath:    "api/vs/facts/999",
			wantResponse:   `{ "error": "fact with ID 999 not found" }`,
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
				if string(body) != tt.wantResponse {
					t.Errorf("got response %s, want %s", string(body), tt.wantResponse)
					return
				}
			}

			_ = resp.Body.Close()
		})
	}

	t.Log("create test")
	t.Fail()
}

func startPublicAPIServer(apiPort string) error {
	go server.Run()

	timeout := time.After(20 * time.Second)
	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-timeout:
			return errors.New("timed out while waiting for public api server to return http status 200 at /Health")
		case <-tick:
			req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%s/Health", apiPort), nil)
			resp, _ := http.DefaultClient.Do(req)
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			_ = resp.Body.Close()
		}
	}
}
