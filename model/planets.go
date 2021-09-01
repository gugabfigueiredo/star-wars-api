package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Planet struct {
	ID primitive.ObjectID 	`bson:"_id,omitempty"`
	Name string				`bson:"name,omitempty"`
	Weather string			`bson:"weather,omitempty"`
	Terrain	string			`bson:"terrain,omitempty"`
	Refs int				`bson:"references,omitempty"`
}