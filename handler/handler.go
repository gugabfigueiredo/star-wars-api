package handler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/model"
	"github.com/gugabfigueiredo/star-wars-api/service"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
)

type IHandler interface {
	service.IService
	FindAllPlanets(w http.ResponseWriter, _ *http.Request)
	FindPlanet(w http.ResponseWriter, r *http.Request)
	FindPlanetByID(w http.ResponseWriter, r *http.Request)
	CreatePlanets(w http.ResponseWriter, r *http.Request)
	SetUpdatedPlanetRefs(w http.ResponseWriter, r *http.Request)
}

type APIHandler struct {
	service.IService
	Logger *log.Logger
}

func (h *APIHandler) FindAllPlanets(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h.Logger.I("Request all planets")

	planets, err := h.GetAllPlanets()
	if err != nil {
		h.Logger.E("Failed to request for all planets", "err", err)
		http.Error(w, "Failed to request for all planets", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&planets); err != nil {
		h.Logger.E("Error on marshal all planets", "err", err)
		http.Error(w, "Error on marshal all planets", http.StatusInternalServerError)
		return
	}
}

func (h *APIHandler) FindPlanetByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := chi.URLParam(r, "name")

	logger := h.Logger.C("name", name)
	logger.I("Request planet by name", "name", name)

	var planet model.Planet
	if err := h.GetPlanet(bson.M{"name": name}, &planet); err != nil {
		logger.E("Error on calling db for planet by name", "err", err)
		http.Error(w, "Error on calling db for planet by name", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(planet); err != nil {
		logger.E("Error on marshal planet by name", "err", err)
		http.Error(w, "Error on marshal planet by name", http.StatusInternalServerError)
		return
	}
}

func (h *APIHandler) FindPlanetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ID := chi.URLParam(r, "planetID")

	logger := h.Logger.C("ID", ID)
	logger.I("Request planet by id", "ID", ID)

	var planet model.Planet
	if err := h.GetPlanet(bson.M{"_id": ID}, &planet);err != nil {
		h.Logger.E("Error on calling db for planet by id", "err", err, "_id", ID)
		http.Error(w, "Error on calling db for planet by id", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(planet); err != nil {
		h.Logger.E("Error on marshal planet by ID", "err", err, "_id", ID, "planet", planet)
		http.Error(w, "Error on marshal planet by ID", http.StatusInternalServerError)
		return
	}
}

func (h *APIHandler) CreatePlanets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h.Logger.I("Landing creation request")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.Logger.E("Error on get body data for planets creation ", "err", err)
		http.Error(w, "Error on get body data for planets creation", http.StatusInternalServerError)
		return
	}

	var planets []*model.Planet
	if err := json.Unmarshal(body, &planets); err != nil {
		h.Logger.E("Error on unmarshal planets payload for creation", "err", err, "planets", planets)
		http.Error(w, "Error on unmarshal planets payload for creation", http.StatusInternalServerError)
	}

	if res, err := h.InsertPlanets(planets); err != nil {
		h.Logger.E("Error on insert planets into database", "err", err, "res", res)
		http.Error(w, "error on insert planets into database", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode([]byte(`{"status":"success"}`)); err != nil {
		h.Logger.E("Error on writing to output stream", "err", err)
		http.Error(w, "Error on writing to output stream", http.StatusInternalServerError)
		return
	}
	return
}

func (h *APIHandler) SetUpdatedPlanetRefs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h.Logger.I("Landing creation request")

	if err := h.UpdatePlanetRefs(); err != nil {
		h.Logger.E("failed to update planet refs by request", "err", err)
		http.Error(w, "failed to update planet refs by request", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode([]byte(`{"status":"success"}`)); err != nil {
		h.Logger.E("Error on writing to output stream", "err", err)
		http.Error(w, "Error on writing to output stream", http.StatusInternalServerError)
		return
	}
	return
}