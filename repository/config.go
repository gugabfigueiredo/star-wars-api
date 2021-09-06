package repository

import (
	"context"
	"fmt"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Context struct {
	context.Context
}

func (ctx *Context) Decode(value string) error {
	switch value {
	default:
		ctx.Context = context.TODO()
	}
	return nil
}

// Config - Configuration for logging
type Config struct {
	Username string `default:"mongo_user"`
	Password string `default:"mongo_pass"`
	Host string `default:"localhost"`
	Port string `default:"27017"`
	AuthDB string `default:"test"`
	Context Context `default:"TODO"`
}

var Repo Repository

func MustInit(config *Config, logger *log.Logger) error {
	Repo = Repository{
		Logger:  logger,
		Context: config.Context,
	}
	// Set client options
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		config.Username, config.Password, config.Host, config.Port, config.AuthDB)
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(Repo.Context, clientOptions)
	if err != nil {
		Repo.Logger.F("failed to connect to database", "err", err)
		return err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		Repo.Logger.F("failed to ping database connection", "err", err)
		return err
	}

	Repo.Client = client
	Repo.Logger.I("connected to database successfully")
	return err
}
