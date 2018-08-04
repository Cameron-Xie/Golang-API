package main

import (
	"gopkg.in/go-playground/validator.v9"
	"github.com/gorilla/mux"
	"go-api/pkg/app"
	"go-api/internal/app/handler/project"
	"go-api/internal/app/middleware"
)

func registerRoutes(db *app.DB) *app.Route {

	routes := &app.Route{Router: mux.NewRouter()}
	routes.Router.Use(middleware.ApiMiddleware)
	v := validator.New()

	// Register Home Handler.
	routes.Get("/projects/", project.List(db))
	routes.Post("/projects/", project.Create(v, db))
	routes.Patch("/projects/{id}", project.Patch(v, db))
	routes.Delete("/projects/{id}", project.Delete(db))

	return routes
}
