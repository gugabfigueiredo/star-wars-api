package test

import (
	"fmt"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/model"
	"github.com/gugabfigueiredo/swapi"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type StringResponse string

type Stub struct {
	Logger *log.Logger

	Planet *model.Planet
	Planets []*model.Planet

	SwapiPlanets []swapi.Planet

	Channel chan bool

	UpdateResult mongo.BulkWriteResult
	InsertResult mongo.InsertManyResult

	CalledWith map[string]interface{}
	RespBody interface{}

	Error error
}

func (s *Stub) GetPlanet(filter interface{}, m *model.Planet) error {
	s.CalledWith = map[string]interface{}{"filter": filter}
	if s.Planet != nil {
		m.Name = s.Planet.Name
		m.Terrain = s.Planet.Terrain
		m.Weather = s.Planet.Weather
		m.Refs = s.Planet.Refs
	}
	return s.Error
}

func (s *Stub) AllPlanets() ([]swapi.Planet, error) {
	return s.SwapiPlanets, s.Error
}

func (s *Stub) GetAllPlanets() ([]*model.Planet, error) {
	return s.Planets, s.Error
}

func (s *Stub) UpdatePlanets(planets []swapi.Planet) (*mongo.BulkWriteResult, error) {
	s.CalledWith = map[string]interface{}{"planets": planets}
	return &s.UpdateResult, s.Error
}

func (s *Stub) InsertPlanets(planets []*model.Planet) (*mongo.InsertManyResult, error) {
	s.CalledWith = map[string]interface{}{"planets": planets}
	return &s.InsertResult, s.Error
}

func (s *Stub) UpdatePlanetRefs() error {
	return s.Error
}
func (s *Stub) SchedulePlanetUpdate(interval time.Duration) chan bool {
	s.CalledWith = map[string]interface{}{"interval": interval}
	return s.Channel
}

func AsString(i interface{}) string {
	return fmt.Sprintf("%+v", i)
}