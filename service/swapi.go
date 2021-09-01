package service

import (
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/swapi"
)

type SwapiService struct {
	*swapi.Client
	Logger *log.Logger
}