package handler

import (
	"encoding/json"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/service"
	"net/http"
)

type HelloHandler struct {
	Service *service.HelloService
	Logger  *log.Logger
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (h *HelloHandler) SayHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qParams := r.URL.Query()
	user := qParams.Get("user")

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
