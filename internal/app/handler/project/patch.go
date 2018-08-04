package project

import (
	"go-api/pkg/app"
	"net/http"
	"gopkg.in/go-playground/validator.v9"
	"go-api/pkg/app/handler"
	"go-api/pkg/functional"
	"go-api/pkg/app/services"
	"go-api/pkg/app/repository"
	"go-api/internal/app/model"
	"github.com/gorilla/mux"
)

func Patch(v *validator.Validate, db *app.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		notFound := handler.JsonErrorHandler(http.StatusNotFound)
		badRequest := handler.JsonErrorHandler(http.StatusBadRequest)

		p, notFoundErr := functional.Compose(
			repository.Find(db.Conn.Where("id = ?", vars["id"]))(notFound))(&model.Project{})

		if notFoundErr != nil {
			handler.JsonResponse(w)(handler.ErrorResponse(notFoundErr))
		} else {
			updated, err := functional.Compose(
				services.Decode(r.Body)(badRequest),
				services.Validate(v)(badRequest),
				repository.Update(db.Conn, p)(badRequest))(&model.Project{})

			handler.JsonResponse(w)(handler.GetResponse(updated, err, http.StatusAccepted))
		}
	}
}
