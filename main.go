package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gugabfigueiredo/star-wars-api/env"
	"github.com/gugabfigueiredo/star-wars-api/handler"
	"github.com/gugabfigueiredo/star-wars-api/log"
	"github.com/gugabfigueiredo/star-wars-api/repository"
	"github.com/gugabfigueiredo/star-wars-api/service"
	"github.com/gugabfigueiredo/swapi"
	"github.com/kelseyhightower/envconfig"
	"net/http"
	"os"
	"time"
)

var Logger *log.Logger

func init() {

	envconfig.MustProcess("swapi", env.Settings)

	Logger = log.New(env.Settings.Log)

	name, _ := os.Hostname()
	Logger = Logger.C("host", name)

	if err := repository.MustInit(env.Settings.Database, Logger); err != nil {
		Logger.F("failed to initialize database connection", "err", err, "settings", env.Settings.Database)
	}
}

func main() {

	// Services

	helloService := &service.HelloService{
		Logger: Logger,
	}

	apiService := &service.APIService{
		Repository:  &repository.Repo,
		SwapiClient: swapi.DefaultClient,
		Logger:      Logger,
	}

	// Handlers

	helloHandler := &handler.HelloHandler{
		Logger: Logger,
		Service: helloService,
	}

	apiHandler := &handler.APIHandler{
		Logger:     Logger,
		APIService: apiService,
	}


	r := chi.NewRouter()

	r.Get("/health", helloHandler.SayHello)

	r.Route("/planets", func(r chi.Router) {
		r.Get("/", apiHandler.FindAllPlanets)
		r.Get("/name/{name:[a-z0-9_]+}", apiHandler.FindPlanetByName)
		r.Get("/id/{planetID:[0-9]+}", apiHandler.FindPlanetByID)
		r.Get("/update", apiHandler.UpdatePlanetRefs)

		r.Post("/create", apiHandler.CreatePlanets)
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

	// update planet movie refs
	schedule := apiService.SchedulePlanetUpdate(env.Settings.Server.UpdateRefsTimeout)

	if err := server.ListenAndServe(); err != nil {
		schedule <- false
		close(schedule)
		Logger.F("listen and serve died", "err", err)
	}
}