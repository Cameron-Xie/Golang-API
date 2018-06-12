package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Router *mux.Router
}

func (route *Route) Get(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	route.Router.HandleFunc(path, handler).Methods("GET")
}
