package service

import (
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/model"
)

type APIService struct {
	Logger *log.Logger
	*repository.APIRepo
}

func (service *APIService) GetAllPlanets() ([]*model.Planet, error) {

}
