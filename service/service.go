package service

import (
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/repository"
	"github.com/gugabfigueiredo/swapi"
	"time"
)

type APIService struct {
	*repository.Repository
	SwapiClient *swapi.Client
	Logger      *log.Logger
}

func (api *APIService) UpdatePlanetRefs() error {
	// get all planets from swapi
	planets, err := api.SwapiClient.AllPlanets()
	if err != nil {
		api.Logger.E("failed to query swapi for planet data", "err", err)
		return err
	}
	// update planets
	res, err := api.UpdatePlanets(planets)
	if err != nil {
		api.Logger.E("failed to write planets to database", "err", err, "result", res)
		return err
	}
	return nil
}

func (api *APIService) SchedulePlanetUpdate(interval time.Duration) chan bool {
	// start ticker to update database
	ticker := time.NewTicker(interval)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <- ticker.C:
				if err := api.UpdatePlanetRefs(); err != nil {
					api.Logger.E("failed to update planet references", "err", err)
				}
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}
