package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

const (
	defaultPage int    = 1
	idKey       ctxKey = "id"
)

type ctxKey string

type Collection struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

type listParam struct {
	Page   int
	Offset int
	Limit  int
}

type StoreService interface {
	Store(r io.ReadCloser) (interface{}, error)
	Delete(uuid.UUID) error
}

type UpdateService interface {
	Update(r io.ReadCloser, id uuid.UUID) error
}

type ReadService interface {
	List(offset, limit int) (*Collection, error)
	Read(uuid.UUID) (interface{}, error)
}

func NewEndpoint(
	storeSvc StoreService,
	updateSvc UpdateService,
	readSvc ReadService,
	render Render,
	limit int,
) http.Handler {
	r := chi.NewRouter()

	r.Get("/", list(readSvc, render, limit))
	r.Post("/", store(storeSvc, render))

	r.Route(fmt.Sprintf("/{%v}", idKey), func(r chi.Router) {
		r.Use(setUUID(idKey, render))
		r.Get("/", read(readSvc, render))
		r.Patch("/", update(updateSvc, render))
		r.Delete("/", delete(storeSvc, render))
	})

	return r
}

func list(s ReadService, render Render, limit int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := getListParam(r.URL.Query(), defaultPage, limit)

		coll, err := s.List(p.Offset, p.Limit)
		if err != nil {
			render.Error(w, err)
			return
		}

		render.Success(w, toCollectionResp(r.URL, p.Page, p.Limit, defaultPage, coll), http.StatusOK)
	}
}

func store(s StoreService, render Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o, err := s.Store(r.Body)

		if err != nil {
			render.Error(w, err)
			return
		}

		render.Success(w, o, http.StatusCreated)
	}
}

func read(s ReadService, render Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		i, _ := ctx.Value(idKey).(uuid.UUID)

		entity, err := s.Read(i)
		if err != nil {
			render.Error(w, err)
			return
		}

		render.Success(w, entity, http.StatusOK)
	}
}

func update(s UpdateService, render Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		i, _ := r.Context().Value(idKey).(uuid.UUID)

		if err := s.Update(r.Body, i); err != nil {
			render.Error(w, err)
			return
		}

		render.Success(w, nil, http.StatusNoContent)
	}
}

func delete(s StoreService, render Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		i, _ := r.Context().Value(idKey).(uuid.UUID)

		if err := s.Delete(i); err != nil {
			render.Error(w, err)
			return
		}

		render.Success(w, nil, http.StatusNoContent)
	}
}

func setUUID(key ctxKey, render Render) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			i, ok := isUUID(chi.URLParam(r, string(key)), 4)

			if !ok {
				render.Error(w, NotFoundError())
				return
			}

			ctx := context.WithValue(r.Context(), key, *i)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
