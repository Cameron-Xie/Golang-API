package project

import (
	"net/http"
	"gopkg.in/go-playground/validator.v9"
	"go-api/internal/app/model"
	"go-api/pkg/app/handler"
	"go-api/pkg/functional"
	"go-api/pkg/app"
	"go-api/pkg/app/services"
	"go-api/pkg/app/repository"
)

func Create(v *validator.Validate, db *app.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		badRequest := handler.JsonErrorHandler(http.StatusBadRequest)

		p, err := functional.Compose(
			services.Decode(r.Body)(badRequest),
			services.Validate(v)(badRequest),
			repository.Create(db.Conn)(badRequest))(&model.Project{})

		handler.JsonResponse(w)(handler.GetResponse(p, err, http.StatusCreated))
	}
}
