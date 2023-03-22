package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cafo13/animal-facts/backend/common/client/facts"
	"github.com/stretchr/testify/require"
)

func authorizationBearer(token string) func(context.Context, *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		return nil
	}
}

type FactsHTTPClient struct {
	client *facts.ClientWithResponses
}

func NewFactsHTTPClient(t *testing.T, token string) FactsHTTPClient {
	addr := os.Getenv("FACTS_HTTP_ADDR")
	fmt.Println("Trying facts http:", addr)
	ok := WaitForPort(addr)
	require.True(t, ok, "Facts HTTP timed out")

	url := fmt.Sprintf("http://%v/api", addr)

	client, err := facts.NewClientWithResponses(
		url,
		facts.WithRequestEditorFn(authorizationBearer(token)),
	)
	require.NoError(t, err)

	return FactsHTTPClient{
		client: client,
	}
}

func (c FactsHTTPClient) CreateFact(t *testing.T, note string, hour time.Time) string {
	response, err := c.client.CreateFactWithResponse(context.Background(), facts.CreateFactJSONRequestBody{
		Notes: note,
		Time:  hour,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode())

	contentLocation := response.HTTPResponse.Header.Get("content-location")

	return lastPathElement(contentLocation)
}

func (c FactsHTTPClient) CreateFactShouldFail(t *testing.T, note string, hour time.Time) {
	response, err := c.client.CreateFact(context.Background(), facts.CreateFactJSONRequestBody{
		Notes: note,
		Time:  hour,
	})
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())

	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func (c FactsHTTPClient) GetFacts(t *testing.T) facts.Facts {
	response, err := c.client.GetFactsWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func (c FactsHTTPClient) CancelFact(t *testing.T, factUUID string, expectedStatusCode int) {
	response, err := c.client.CancelFact(context.Background(), factUUID)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())

	require.Equal(t, expectedStatusCode, response.StatusCode)
}
