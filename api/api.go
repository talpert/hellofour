package api

import (
	"net/http"
	"os"

	hh "github.com/InVisionApp/go-health/handlers"
	"github.com/InVisionApp/rye"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/handlers"
	"github.com/talpert/hellofour/config"
	"github.com/talpert/hellofour/deps"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("pkg", "api")
}

type API struct {
	Config  *config.Config
	Version string
	Deps    *deps.Dependencies
}

type APIResponseJSON struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Values  map[string]string `json:"values,omitempty"`
	Errors  string            `json:"errors,omitempty"`
}

func New(cfg *config.Config, d *deps.Dependencies, version string) *API {
	return &API{
		Config:  cfg,
		Version: version,
		Deps:    d,
	}
}

func (a *API) Run() error {
	llog := log.WithField("method", "Run")
	llog.Infof("Starting API server...")

	routes := mux.NewRouter().StrictSlash(true)

	/***************
	Basic handlers
	***************/

	routes.Handle(
		"/", http.HandlerFunc(a.homeHandler),
	).Methods("GET")

	routes.Handle(
		"/version", http.HandlerFunc(a.versionHandler),
	).Methods("GET")

	healthHandler := hh.NewJSONHandlerFunc(a.Deps.Health, map[string]interface{}{
		"version": a.Version,
	})

	routes.Handle(
		"/healthcheck", healthHandler,
	).Methods("GET")

	/*************
	v1 endpoints
	*************/

	routes.Handle(a.setupHandler(
		"/v1/heroku/resources", []rye.Handler{
			rye.NewMiddlewareAuth(rye.NewBasicAuthFunc(map[string]string{
				"user1":     "my_password",
				"hellofour": "80e6fa24349745346f434e630be2e456",
			})),
			a.resourceHandler,
		})).Methods("POST")

	llog.Infof("API server running on :%v", a.Config.ListenAddress)

	return http.ListenAndServe(":"+a.Config.ListenAddress, routes)
}

func (a *API) setupHandler(path string, ryeStack []rye.Handler) (string, http.Handler) {
	return path, handlers.LoggingHandler(os.Stdout, a.Deps.MWHandler.Handle(ryeStack))
}
