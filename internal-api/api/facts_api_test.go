//go1:build integration

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/cafo13/animal-facts/internal-api/handler"
	"github.com/cafo13/animal-facts/pkg/repository"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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

func Test_RunApiIntegrationTests_create_update_delete(t *testing.T) {
	factsApi, err := getFactsApi()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("create, update and delete a fact", func(t *testing.T) {
		e := echo.New()
		createReqBody := CreateUpdateFact{
			Fact:     "Some new animal fact.",
			Source:   "some-source.com/animalfact/23",
			Approved: true,
		}
		createJsonBody, err := json.Marshal(createReqBody)
		if err != nil {
			t.Errorf("Test_RunApiIntegrationTests_create_update_delete() error at marshalling create request body = %v", err)
			return
		}
		createReq := httptest.NewRequest(http.MethodPost, "/facts", bytes.NewBuffer(createJsonBody))
		createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		createRec := httptest.NewRecorder()
		createCtx := e.NewContext(createReq, createRec)

		err = factsApi.createFact(createCtx)
		if err != nil {
			t.Errorf("Test_RunApiIntegrationTests_create_update_delete() unepected error in create fact handler = %v", err)
			return
		}

		if createRec.Code != http.StatusCreated {
			t.Errorf("Test_RunApiIntegrationTests_create_update_delete() expected status code 201 for create response, got %d", createRec.Code)
			return
		}

		var createRespBody CreateFactResult
		createRespBodyString := strings.TrimSuffix(createRec.Body.String(), "\n")
		err = json.Unmarshal([]byte(createRespBodyString), &createRespBody)
		if err != nil {
			t.Errorf("Test_RunApiIntegrationTests_create_update_delete() unepected error at unmarshalling create response json = %v", err)
			return
		}
		if createRespBody.Id == "" {
			t.Errorf("Test_RunApiIntegrationTests_create_update_delete() expected fact id to be set in create response, got empty string, full resp string = %s", createRespBodyString)
			return
		}

		updateReqBody := CreateUpdateFact{
			Fact:     "Some updated animal fact.",
			Source:   "some-source.com/animalfact/23",
			Approved: true,
		}
		updateJsonBody, err := json.Marshal(updateReqBody)
		if err != nil {
			t.Errorf("Test_RunApiIntegrationTests_create_update_delete() error at marshalling update request body = %v", err)
			return
		}
		updateReq := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/facts/%s", createRespBody.Id), bytes.NewBuffer(updateJsonBody))
		updateReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		updateRec := httptest.NewRecorder()
		updateCtx := e.NewContext(updateReq, updateRec)
		updateCtx.SetParamNames("id")
		updateCtx.SetParamValues(createRespBody.Id)

		err = factsApi.updateFact(updateCtx)
		if err != nil {
			t.Errorf("Test_RunApiIntegrationTests_create_update_delete() unepected error at updating fact = %v", err)
			return
		}

		if updateRec.Code != http.StatusOK {
			t.Errorf("Test_RunApiIntegrationTests_create_update_delete() expected status code 200 for update response, got %d", updateRec.Code)
			return
		}

		deleteReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/facts/%s", createRespBody.Id), nil)
		deleteRec := httptest.NewRecorder()
		deleteCtx := e.NewContext(deleteReq, deleteRec)
		deleteCtx.SetParamNames("id")
		deleteCtx.SetParamValues(createRespBody.Id)

		err = factsApi.deleteFact(deleteCtx)
		if err != nil {
			t.Errorf("Test_RunApiIntegrationTests_create_update_delete() unepected error at deleting fact = %v", err)
			return
		}
	})
}
