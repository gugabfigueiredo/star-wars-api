package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/model"
	"github.com/gugabfigueiredo/star-wars-api/test"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIHandler_FindPlanets(t *testing.T) {
	tests := []struct{
		name               	string
		stub               	*test.Stub
		endpoint			string
		pathParam			string
		expectedStatusCode 	int
		expectedBody       	interface{}
		expectedListBody   	[]map[string]interface{}
		expectedErr        	error
		expectedCalledWith 	map[string]interface{}
		expectedContentType	string
	}{
		{
			name: "find planet by name",
			stub: &test.Stub{Planet: &model.Planet{Name: "Planet", Weather: "nice", Terrain: "rocky", Refs: 0}, RespBody: &model.Planet{}},
			endpoint: "name",
			pathParam: "name",
			expectedBody: &model.Planet{Name: "Planet", Weather: "nice", Terrain: "rocky", Refs: 0},
			expectedStatusCode: http.StatusOK,
			expectedContentType: "application/json",
			expectedCalledWith: map[string]interface{}{"filter": bson.M{"name": "name"}},
		},
		{
			name: "find planet by id",
			stub: &test.Stub{Planet: &model.Planet{Name: "Planet", Weather: "nice", Terrain: "rocky", Refs: 0}, RespBody: &model.Planet{}},
			endpoint: "id",
			pathParam: "1234",
			expectedBody: &model.Planet{Name: "Planet", Weather: "nice", Terrain: "rocky", Refs: 0},
			expectedStatusCode: http.StatusOK,
			expectedContentType: "application/json",
			expectedCalledWith: map[string]interface{}{"filter": bson.M{"_id": "1234"}},
		},
		{
			name: "no planet match by name",
			stub: &test.Stub{},
			endpoint: "name",
			pathParam: "name",
			expectedContentType: "application/json",
			expectedStatusCode: http.StatusOK,
			expectedCalledWith: map[string]interface{}{"filter": bson.M{"name": "name"}},
		},
		{
			name: "no planet match by id",
			stub: &test.Stub{},
			endpoint: "id",
			pathParam: "1234",
			expectedContentType: "application/json",
			expectedStatusCode: http.StatusOK,
			expectedCalledWith: map[string]interface{}{"filter": bson.M{"_id": "1234"}},
		},
		{
			name: "fail to query for planet by name",
			stub: &test.Stub{Error: errors.New("fail to query for planet by name")},
			endpoint: "name",
			pathParam: "name",
			expectedStatusCode: http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
			expectedCalledWith: map[string]interface{}{"filter": bson.M{"name": "name"}},
		},
		{
			name: "fail to query for planet by id",
			stub: &test.Stub{Error: errors.New("fail to query for planet by id")},
			endpoint: "id",
			pathParam: "1234",
			expectedStatusCode: http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
			expectedCalledWith: map[string]interface{}{"filter": bson.M{"_id": "1234"}},
		},
		{
			name: "get all planets",
			stub: &test.Stub{
				Planets: []*model.Planet{
					{Name: "Planet1", Weather: "nice", Terrain: "rocky", Refs: 1},
					{Name: "Planet2", Weather: "warm", Terrain: "icy", Refs: 2},
					{Name: "Planet3", Weather: "cold", Terrain: "plains", Refs: 3},
					{Name: "Planet4", Weather: "nice", Terrain: "forest", Refs: 4},
				},
				RespBody: &[]*model.Planet{},
			},
			endpoint: "planets",
			expectedBody: &[]*model.Planet{
				{Name: "Planet1", Weather: "nice", Terrain: "rocky", Refs: 1},
				{Name: "Planet2", Weather: "warm", Terrain: "icy", Refs: 2},
				{Name: "Planet3", Weather: "cold", Terrain: "plains", Refs: 3},
				{Name: "Planet4", Weather: "nice", Terrain: "forest", Refs: 4},
			},
			expectedStatusCode: http.StatusOK,
			expectedContentType: "application/json",
		},
		{
			name: "fail to get all planets",
			stub: &test.Stub{Error: errors.New("failed to query db for all planets")},
			endpoint: "planets",
			expectedStatusCode: http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
		},
	}

	logger := log.New(&log.Config{
		Context:               "sw-api-test",
		ConsoleLoggingEnabled: false,
		EncodeLogsAsJson:      true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stub.Logger = logger
			h := &APIHandler{
				IService: tt.stub,
				Logger:   logger,
			}

			router := chi.NewRouter()
			router.Get("/name/{name:[a-z0-9_]+}", h.FindPlanetByName)
			router.Get("/id/{planetID:[0-9]+}", h.FindPlanetByID)
			router.Get("/planets", h.FindAllPlanets)
			mockServer := httptest.NewServer(router)
			defer mockServer.Close()

			endpoint := fmt.Sprintf("%s/%s", mockServer.URL, tt.endpoint)
			if tt.pathParam != "" {
				endpoint = fmt.Sprintf("%s/%s", endpoint, tt.pathParam)
			}
			r, err := http.NewRequest(http.MethodGet, endpoint, nil)
			if err != nil {
				t.Fatalf("could not create request. err %+v\n", err)
			}

			resp, err := http.DefaultClient.Do(r)
			assert.Equal(t, test.AsString(tt.expectedErr), test.AsString(err))
			if err == nil {
				defer resp.Body.Close()
				assert.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"))
				assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
				assert.Equal(t, test.AsString(tt.expectedCalledWith), test.AsString(tt.stub.CalledWith))

				if tt.expectedBody != nil {
					if err := json.NewDecoder(resp.Body).Decode(tt.stub.RespBody); err != nil {
						t.Fatalf("could not read the response body. err %+v\n", err)
					}

					assert.Equal(t, tt.expectedBody, tt.stub.RespBody)
				}
			}
		})
	}
}

func TestAPIHandler_CreatePlanets(t *testing.T) {

	tests := []struct{
		name               	string
		stub               	*test.Stub
		requestBody        	map[string]interface{}
		endpoint			string
		expectedStatusCode 	int
		expectedBody       	interface{}
		expectedErr        	error
		expectedCalledWith 	map[string]interface{}
		expectedContentType	string
	}{
		{
			name: "create one planet",
			requestBody: map[string]interface{}{"planets": []*model.Planet{
				{Name: "NewPlanet", Weather: "nice", Terrain: "slimy", Refs: 0},
			}},
			expectedStatusCode: http.StatusOK,
		},
		{name: "create many planets"},
		{name: "fail to create"},
	}

	logger := log.New(&log.Config{
		Context:               "sw-api-test",
		ConsoleLoggingEnabled: false,
		EncodeLogsAsJson:      true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stub.Logger = logger
			h := &APIHandler{
				IService: tt.stub,
				Logger:   logger,
			}

			router := chi.NewRouter()
			router.Get("/create", h.CreatePlanets)
			mockServer := httptest.NewServer(router)
			defer mockServer.Close()

			endpoint := fmt.Sprintf("%s/%s", mockServer.URL, tt.endpoint)
			if tt.pathParam != "" {
				endpoint = fmt.Sprintf("%s/%s", endpoint, tt.pathParam)
			}
			r, err := http.NewRequest(http.MethodGet, endpoint, nil)
			if err != nil {
				t.Fatalf("could not create request. err %+v\n", err)
			}

			resp, err := http.DefaultClient.Do(r)
			assert.Equal(t, test.AsString(tt.expectedErr), test.AsString(err))
			if err == nil {
				defer resp.Body.Close()
				assert.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"))
				assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
				assert.Equal(t, test.AsString(tt.expectedCalledWith), test.AsString(tt.stub.CalledWith))

				if tt.expectedBody != nil {
					if err := json.NewDecoder(resp.Body).Decode(tt.stub.RespBody); err != nil {
						t.Fatalf("could not read the response body. err %+v\n", err)
					}

					assert.Equal(t, tt.expectedBody, tt.stub.RespBody)
				}
			}
		})
	}
}

func TestAPIHandler_UpdatePlanetRefs(t *testing.T) {

	tests := []struct{
		name               	string
		stub               	*test.Stub
		requestBody        	map[string]interface{}
		endpoint			string
		pathParam			string
		expectedStatusCode 	int
		expectedBody       	interface{}
		expectedErr        	error
		expectedCalledWith 	map[string]interface{}
		expectedContentType	string
	}{
		{name: "update planet refs"},
		{name: "fail to update planet refs"},
	}
}
