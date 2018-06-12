package main

import (
	"app"
	"app/handler"

	"github.com/gorilla/mux"
)

func registerRoutes() *app.Route {

	routes := &app.Route{Router: mux.NewRouter()}

	// Register Home Handler.
	home := &handler.HomeHandler{}
	routes.Get("/", home.Welcome)

	return routes
}
