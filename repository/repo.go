package repository

import (
	"context"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/model"
	"github.com/gugabfigueiredo/swapi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	*mongo.Client
	Context context.Context
	Logger *log.Logger
}

func (r *Repository) Planets() *mongo.Collection {
	return r.Database("sw-api").Collection("planets")
}

func (r *Repository) GetPlanet(filter interface{}, model interface{}) error {
	return r.Planets().FindOne(r.Context, filter).Decode(model)
}

func (r *Repository) GetAllPlanets() ([]*model.Planet, error) {

	cur, err := r.Planets().
		Find(r.Context, bson.D{})
	if err != nil {
		r.Logger.E("failed to query for planets", "err", err)
		return nil, err
	}

	var results []*model.Planet
	for cur.Next(r.Context) {

		var planet *model.Planet
		if err := cur.Decode(&planet); err != nil {
			r.Logger.E("failed to decode planet", "err", err, "planet", planet)
			return nil, err
		}

		results = append(results, planet)
	}

	if err := cur.Err(); err != nil {
		r.Logger.E("error at the end of cursor", "err", err)
		return nil, err
	}

	return results, nil
}

func (r *Repository) UpdatePlanets(planets []swapi.Planet) (*mongo.BulkWriteResult, error) {
	var writes []mongo.WriteModel
	for _, planet := range planets {
		writes = append(writes, model.WritePlanetModel(planet))
	}
	return r.Planets().BulkWrite(r.Context, writes, options.BulkWrite().SetOrdered(false))
}

func (r *Repository) InsertPlanets(planets []*model.Planet) (*mongo.InsertManyResult ,error) {

	var docs []interface{}
	for _, planet := range planets {
		data, err := bson.Marshal(planet)
		if err != nil {
			r.Logger.E("failed to marshal planet", "err", err, "planet", planet)
			return nil, err
		}

		var doc bson.D
		if err := bson.Unmarshal(data, &doc); err != nil {
			r.Logger.E("failed to unmarshal data into planet bson.D", "err", err, "data", data)
		}
		docs = append(docs, data)
	}

	return r.Planets().InsertMany(r.Context, docs, options.InsertMany().SetOrdered(false))
}