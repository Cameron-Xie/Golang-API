package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(v int, endpoints map[string]http.Handler, l middleware.LogFormatter, c *cors.Cors) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger(l))
	r.Use(middleware.Recoverer)

	if c != nil {
		r.Use(c.Handler)
	}

	r.Route(fmt.Sprintf("/v%v", v), func(r chi.Router) {
		for p, e := range endpoints {
			r.Mount(p, e)
		}
	})

	return r
}

func SetCORS(origins []string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
