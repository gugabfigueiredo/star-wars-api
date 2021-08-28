package handler

import (
	"encoding/json"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/service"
	"net/http"
)

type APIHandler struct {
	Logger *log.Logger
	*service.APIService
}

func New(logger *log.Logger, service *service.APIService) *APIHandler {
	return &APIHandler{
		Logger: logger,
		APIService: service,
	}
}

func (h *APIHandler) FindAllPlanets(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	h.Logger.I("Request all planets")

	planets, err := h.GetAllPlanets()
	if err != nil {
		h.Logger.E("Error on calling db for all planets", "err", err)
		http.Error(w, "Error on calling db for all planets", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(planets); err != nil {
		h.Logger.E("Error on marshal all planets", "err", err)
		http.Error(w, "Error on marshal all planets", http.StatusInternalServerError)
		return
	}
}