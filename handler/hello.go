package handler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/service"
	"net/http"
)

type HelloHandler struct {
	Logger  *log.Logger
	Service *service.HelloService
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (h *HelloHandler) SayHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Read request params
	user := chi.URLParam(r, "user")

	logger := h.Logger.C("user", user)

	// Call service
	message, err := h.Service.SayHello(user)

	response := Response{
		Status:  "OK",
		Message: message,
	}

	if err != nil {
		response.Status = "ERROR"
		response.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Write response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.E("error on json encoding", "err", err)
		http.Error(w, "error while writing response", http.StatusInternalServerError)
		return
	}
}
