package model

import (
	"github.com/gugabfigueiredo/swapi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SwapiWritePlanetModel(planet *swapi.Planet) mongo.WriteModel {
	model := mongo.NewUpdateOneModel()
	model.SetFilter(bson.M{"name": planet.Name})
	model.SetUpdate(bson.M{"$set": bson.M{
		"name":       planet.Name,
		"weather":    planet.Climate,
		"terrain":    planet.Terrain,
		"references": len(planet.FilmURLs),
	}})
	model.SetUpsert(true)
	return model
}

func WritePlanetModel(planet *Planet) mongo.WriteModel {
	model := mongo.NewUpdateOneModel()
	model.SetFilter(bson.M{"name": planet.Name})
	model.SetUpdate(bson.M{"$set": bson.M{
		"name":       planet.Name,
		"weather":    planet.Climate,
		"terrain":    planet.Terrain,
		"references": planet.Refs,
	}})
	return model
}