package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Router *mux.Router
}

func (r *Route) Get(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.Router.HandleFunc(path, handler).Methods("GET")
}

func (r *Route) Post(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.Router.HandleFunc(path, handler).Methods("POST")
}

func (r *Route) Patch(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.Router.HandleFunc(path, handler).Methods("PATCH")
}

func (r *Route) Delete(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.Router.HandleFunc(path, handler).Methods("DELETE")
}