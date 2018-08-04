package project

import (
	"net/http"
	"go-api/pkg/app"
	"go-api/internal/app/model"
	"go-api/pkg/app/handler"
	"go-api/pkg/functional"
	"github.com/gorilla/mux"
	"go-api/pkg/app/repository"
)

func Delete(db *app.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		notFound := handler.JsonErrorHandler(http.StatusNotFound)

		p, err := functional.Compose(
			repository.Find(db.Conn.Where("id = ?", vars["id"]))(notFound),
			repository.Delete(db.Conn)(notFound))(&model.Project{})

		handler.JsonResponse(w)(handler.GetResponse(p, err, http.StatusNoContent))
	}
}
