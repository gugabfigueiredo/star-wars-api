package test

import (
	"fmt"
	"github.com/gugabfigueiredo/star-wars-api/model"
	"github.com/gugabfigueiredo/swapi"
	"go.mongodb.org/mongo-driver/mongo"
)

type StringResponse string

type Stub struct {

	Planet *model.Planet
	Planets []*model.Planet

	SwapiPlanets []swapi.Planet
	SwapiUpdates int

	Channel chan bool

	InsertResult mongo.InsertManyResult
	UpdateResult  mongo.BulkWriteResult
	DeleteResult mongo.DeleteResult

	CalledWith map[string]interface{}
	RespBody interface{}

	Error error
}

func (s *Stub) GetPlanet(filter interface{}, m *model.Planet) error {
	s.CalledWith = map[string]interface{}{"filter": filter}
	if s.Planet != nil {
		m.Name = s.Planet.Name
		m.Terrain = s.Planet.Terrain
		m.Climate = s.Planet.Climate
		m.Refs = s.Planet.Refs
	}
	return s.Error
}

func (s *Stub) GetAllPlanets() ([]*model.Planet, error) {
	return s.Planets, s.Error
}

func (s *Stub) UpdateMovieRefs(planets []swapi.Planet) (*mongo.BulkWriteResult, error) {
	s.CalledWith = map[string]interface{}{"planets": planets}
	return &s.UpdateResult, s.Error
}

func (s *Stub) InsertPlanets(planets []model.Planet) (*mongo.InsertManyResult, error) {
	s.CalledWith = map[string]interface{}{"planets": planets}
	return &s.InsertResult, s.Error
}

func (s *Stub) UpdatePlanets(planet []model.Planet) (*mongo.BulkWriteResult, error) {
	s.CalledWith = map[string]interface{}{"planets": planet}
	return &s.UpdateResult, s.Error
}

func (s *Stub) DeletePlanets(planets []model.Planet) (*mongo.DeleteResult, error) {
	s.CalledWith = map[string]interface{}{"planets": planets}
	return &s.DeleteResult, s.Error
}

func (s *Stub) UpdatePlanetRefs() error {
	return s.Error
}

func AsString(i interface{}) string {
	return fmt.Sprintf("%+v", i)
}