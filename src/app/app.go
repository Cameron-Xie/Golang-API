package app

import (
	"log"
	"net/http"
)

type App struct {
	Config *Config
	Route  *Route
}

func (app *App) Run() {
	log.Fatal(http.ListenAndServe(app.Config.Host, app.Route.Router))
}
