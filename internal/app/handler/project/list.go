package project

import (
	"net/http"
	"go-api/pkg/app"
	"go-api/internal/app/model"
	"go-api/pkg/app/handler"
	"go-api/pkg/app/repository"
)

func List(db *app.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		notFound := handler.JsonErrorHandler(http.StatusBadRequest)
		p, err := repository.Find(db.Conn)(notFound)(&[]model.Project{})

		handler.JsonResponse(w)(handler.GetResponse(p, err, http.StatusOK))
	}
}
