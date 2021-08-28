package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gugabfigueiredo/star-wars-api/env"
	"github.com/gugabfigueiredo/star-wars-api/handler"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/service"
	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"os"
	"time"
)

var Logger *log.Logger

func init() {

	envconfig.MustProcess("SWAPI", env.Settings)

	Logger = log.New(env.Settings.Log)

	name, _ := os.Hostname()
	Logger = Logger.C("host", name)

}

func main() {

	helloService := &service.HelloService{
		Logger: Logger,
	}

	hello := handler.HelloHandler{
		Logger: Logger,
		Service: helloService,
	}


	r := chi.NewRouter()

	r.Get("/health", hello.SayHello)

	r.Route("/planets", func(r chi.Router) {
		//r.Get("/", handler.FindAllPlanets)
		//r.Get("/name/{name:[a-z0-9_]+}", handler.FindPlanetByName)
		//r.Get("/id/{carrier_id:[0-9]+}", handler.FindPlanetByID)
		//
		//r.Post("/create", handler.CreatePlanets)
	})

	http.Handle("/", r)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", env.Settings.Server.Port),
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	Logger.I("Starting server...", "port", env.Settings.Server.Port)

	if err := server.ListenAndServe(); err != nil {
		Logger.F("listen and serve died", "err", err)
	}
}