package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Cameron-Xie/Golang-API/internal/storage/postgres"
	"github.com/Cameron-Xie/Golang-API/pkg/http/rest"
)

type jsonRender struct {
}

func (r *jsonRender) Success(w http.ResponseWriter, i interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	b, err := toJSON(i)
	if err != nil {
		r.Error(w, err)
		return
	}

	_, _ = w.Write(b)
}

func (r *jsonRender) Error(w http.ResponseWriter, err error) {
	httpErr := toHTTPError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.StatusCode)
	_, _ = w.Write([]byte(httpErr.Error()))
}

func NewJSONRender() rest.Render {
	return new(jsonRender)
}

func toJSON(i interface{}) ([]byte, error) {
	switch b := i.(type) {
	case string:
		return isValidJSON([]byte(b))
	case []byte:
		return isValidJSON(b)
	default:
		return json.Marshal(i)
	}
}

func isValidJSON(i []byte) ([]byte, error) {
	if json.Valid(i) {
		return i, nil
	}

	return nil, errors.New("invalid json")
}

func toHTTPError(err error) *rest.HTTPError {
	switch e := err.(type) {
	case *rest.HTTPError:
		return e
	case *postgres.NotFoundError:
		return &rest.HTTPError{
			Title:      e.Error(),
			StatusCode: http.StatusNotFound,
		}
	default:
		return &rest.HTTPError{
			Title:      e.Error(),
			StatusCode: http.StatusBadGateway,
		}
	}
}
