package handler

import (
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/service"
)

type SwapiHandler struct {
	*service.SwapiService
	Logger *log.Logger
}

