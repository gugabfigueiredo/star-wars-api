package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/model"
	"github.com/gugabfigueiredo/star-wars-api/test"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
			stub: &test.Stub{Planet: &model.Planet{Name: "Planet", Climate: "nice", Terrain: "rocky", Refs: 0}, RespBody: &model.Planet{}},
			endpoint: "name",
			pathParam: "name",
			expectedBody: &model.Planet{Name: "Planet", Climate: "nice", Terrain: "rocky", Refs: 0},
			expectedStatusCode: http.StatusOK,
			expectedContentType: "application/json",
			expectedCalledWith: map[string]interface{}{"filter": bson.M{"name": "name"}},
		},
		{
			name: "find planet by id",
			stub: &test.Stub{Planet: &model.Planet{Name: "Planet", Climate: "nice", Terrain: "rocky", Refs: 0}, RespBody: &model.Planet{}},
			endpoint: "id",
			pathParam: "1234",
			expectedBody: &model.Planet{Name: "Planet", Climate: "nice", Terrain: "rocky", Refs: 0},
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
					{Name: "Planet1", Climate: "nice", Terrain: "rocky", Refs: 1},
					{Name: "Planet2", Climate: "warm", Terrain: "icy", Refs: 2},
					{Name: "Planet3", Climate: "cold", Terrain: "plains", Refs: 3},
					{Name: "Planet4", Climate: "nice", Terrain: "forest", Refs: 4},
				},
				RespBody: &[]*model.Planet{},
			},
			endpoint: "planets",
			expectedBody: &[]*model.Planet{
				{Name: "Planet1", Climate: "nice", Terrain: "rocky", Refs: 1},
				{Name: "Planet2", Climate: "warm", Terrain: "icy", Refs: 2},
				{Name: "Planet3", Climate: "cold", Terrain: "plains", Refs: 3},
				{Name: "Planet4", Climate: "nice", Terrain: "forest", Refs: 4},
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

func TestAPIHandler_CreateUpdateDelete(t *testing.T) {

	tests := []struct{
		name               	string
		stub               	*test.Stub
		requestBody        	interface{}
		endpoint			string
		expectedStatusCode 	int
		expectedBody       	interface{}
		expectedErr        	error
		expectedCalledWith 	map[string]interface{}
		expectedContentType	string
	}{
		{
			name: "create one planet",
			stub: &test.Stub{
				InsertResult: mongo.InsertManyResult{InsertedIDs: []interface{}{0}},
				RespBody: &mongo.InsertManyResult{},
			},
			requestBody: []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			},
			endpoint: "create",
			expectedStatusCode: http.StatusOK,
			expectedContentType: "application/json",
			expectedCalledWith: map[string]interface{}{"planets": []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			}},
			expectedBody: &mongo.InsertManyResult{InsertedIDs: []interface{}{0}},
		},
		{
			name: "create many planets",
			stub: &test.Stub{
				InsertResult: mongo.InsertManyResult{InsertedIDs: []interface{}{0,1,2}},
				RespBody: &mongo.InsertManyResult{},
			},
			requestBody: []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
				{Name: "NewPlanet", Climate: "warm", Terrain: "slimy", Refs: 1},
				{Name: "NewPlanet", Climate: "cold", Terrain: "slimy", Refs: 2},
			},
			endpoint: "create",
			expectedStatusCode: http.StatusOK,
			expectedContentType: "application/json",
			expectedCalledWith: map[string]interface{}{"planets": []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
				{Name: "NewPlanet", Climate: "warm", Terrain: "slimy", Refs: 1},
				{Name: "NewPlanet", Climate: "cold", Terrain: "slimy", Refs: 2},
			}},
			expectedBody: &mongo.InsertManyResult{InsertedIDs: []interface{}{0,1,2}},
		},
		{
			name: "fail to decode request body",
			stub: &test.Stub{RespBody: &mongo.InsertManyResult{}},
			requestBody: []map[string]interface{}{{"name": 1}},
			endpoint: "create",
			expectedStatusCode: http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
		},
		{
			name: "fail to create planet",
			stub: &test.Stub{RespBody: &mongo.InsertManyResult{}, Error: errors.New("failed to insert planets")},
			requestBody: []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			},
			endpoint: "create",
			expectedStatusCode: http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
			expectedCalledWith: map[string]interface{}{"planets": []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			}},
		},
		{
			name: "update one planet",
			stub: &test.Stub{
				UpdateResult: mongo.BulkWriteResult{MatchedCount: 1, ModifiedCount: 1},
				RespBody: &mongo.BulkWriteResult{},
			},
			requestBody: []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			},
			endpoint: "update",
			expectedStatusCode: http.StatusOK,
			expectedContentType: "application/json",
			expectedCalledWith: map[string]interface{}{"planets": []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			}},
			expectedBody: &mongo.BulkWriteResult{MatchedCount: 1, ModifiedCount: 1},
		},
		{
			name: "update many planets",
			stub: &test.Stub{
				UpdateResult: mongo.BulkWriteResult{MatchedCount: 3, ModifiedCount: 3},
				RespBody: &mongo.BulkWriteResult{},
			},
			requestBody: []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
				{Name: "NewPlanet", Climate: "warm", Terrain: "slimy", Refs: 1},
				{Name: "NewPlanet", Climate: "cold", Terrain: "slimy", Refs: 2},
			},
			endpoint: "update",
			expectedStatusCode: http.StatusOK,
			expectedContentType: "application/json",
			expectedCalledWith: map[string]interface{}{"planets": []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
				{Name: "NewPlanet", Climate: "warm", Terrain: "slimy", Refs: 1},
				{Name: "NewPlanet", Climate: "cold", Terrain: "slimy", Refs: 2},
			}},
			expectedBody: &mongo.BulkWriteResult{MatchedCount: 3, ModifiedCount: 3},
		},
		{
			name: "fail to decode request body",
			stub: &test.Stub{RespBody: &mongo.BulkWriteResult{}},
			requestBody: []map[string]interface{}{{"name": 1}},
			endpoint: "update",
			expectedStatusCode: http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
		},
		{
			name: "fail to update planets",
			stub: &test.Stub{RespBody: &mongo.BulkWriteResult{}, Error: errors.New("failed to update planets")},
			requestBody: []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			},
			endpoint: "update",
			expectedStatusCode: http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
			expectedCalledWith: map[string]interface{}{"planets": []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			}},
		},
		{
			name: "delete one planet",
			stub: &test.Stub{
				DeleteResult: mongo.DeleteResult{DeletedCount: 1},
				RespBody: &mongo.DeleteResult{},
			},
			requestBody: []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			},
			endpoint: "delete",
			expectedStatusCode: http.StatusOK,
			expectedContentType: "application/json",
			expectedCalledWith: map[string]interface{}{"planets": []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			}},
			expectedBody: &mongo.DeleteResult{DeletedCount: 1},
		},
		{
			name: "delete many planets",
			stub: &test.Stub{
				DeleteResult: mongo.DeleteResult{DeletedCount: 3},
				RespBody: &mongo.DeleteResult{},
			},
			requestBody: []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
				{Name: "NewPlanet", Climate: "warm", Terrain: "slimy", Refs: 1},
				{Name: "NewPlanet", Climate: "cold", Terrain: "slimy", Refs: 2},
			},
			endpoint: "delete",
			expectedStatusCode: http.StatusOK,
			expectedContentType: "application/json",
			expectedCalledWith: map[string]interface{}{"planets": []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
				{Name: "NewPlanet", Climate: "warm", Terrain: "slimy", Refs: 1},
				{Name: "NewPlanet", Climate: "cold", Terrain: "slimy", Refs: 2},
			}},
			expectedBody: &mongo.DeleteResult{DeletedCount: 3},
		},
		{
			name: "fail to decode request body",
			stub: &test.Stub{RespBody: &mongo.DeleteResult{}},
			requestBody: []map[string]interface{}{{"name": 1}},
			endpoint: "delete",
			expectedStatusCode: http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
		},
		{
			name: "fail to create planet",
			stub: &test.Stub{RespBody: &mongo.DeleteResult{}, Error: errors.New("failed to insert planets")},
			requestBody: []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			},
			endpoint: "delete",
			expectedStatusCode: http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
			expectedCalledWith: map[string]interface{}{"planets": []model.Planet{
				{Name: "NewPlanet", Climate: "nice", Terrain: "slimy", Refs: 0},
			}},
		},
	}

	logger := log.New(&log.Config{
		Context:               "sw-api-test",
		ConsoleLoggingEnabled: false,
		EncodeLogsAsJson:      true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &APIHandler{
				IService: tt.stub,
				Logger:   logger,
			}

			router := chi.NewRouter()
			router.Post("/create", h.CreatePlanets)
			router.Post("/update", h.PlanetUpdate)
			router.Post("/delete", h.RemovePlanets)
			mockServer := httptest.NewServer(router)
			defer mockServer.Close()

			b, err := json.Marshal(map[string]interface{}{})
			if tt.requestBody != nil {
				b, err = json.Marshal(tt.requestBody)
			}
			if err != nil {
				t.Fatalf("could not marshal request body")
			}

			endpoint := fmt.Sprintf("%s/%s", mockServer.URL, tt.endpoint)
			r, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
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

					assert.Equal(t, test.AsString(tt.expectedBody), test.AsString(tt.stub.RespBody))
				}
			}
		})
	}
}

func TestAPIHandler_UpdatePlanetRefs(t *testing.T) {

	tests := []struct{
		name               	string
		stub               	*test.Stub
		expectedStatusCode 	int
		expectedErr        	error
		expectedContentType	string
	}{
		{
			name: "update planet refs",
			stub: &test.Stub{},
			expectedStatusCode: http.StatusOK,
			expectedContentType: "application/json",
		},
		{
			name: "fail to get updated refs",
			stub: &test.Stub{Error: errors.New("failed to get updated refs")},
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
			h := &APIHandler{
				IService: tt.stub,
				Logger:   logger,
			}

			router := chi.NewRouter()
			router.Get("/", h.SetMovieRefs)
			mockServer := httptest.NewServer(router)
			defer mockServer.Close()

			r, err := http.NewRequest(http.MethodGet, mockServer.URL, nil)
			if err != nil {
				t.Fatalf("could not create request. err %+v\n", err)
			}

			resp, err := http.DefaultClient.Do(r)
			assert.Equal(t, test.AsString(tt.expectedErr), test.AsString(err))
			if err == nil {
				defer resp.Body.Close()
				assert.Equal(t, tt.expectedContentType, resp.Header.Get("Content-Type"))
				assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}