package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Planet struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name,omitempty"`
	Climate string             `bson:"weather,omitempty"`
	Terrain string             `bson:"terrain,omitempty"`
	Refs    int                `bson:"references,omitempty"`
}