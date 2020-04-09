package rest

import (
	"net/http"
)

type Render interface {
	Success(w http.ResponseWriter, i interface{}, code int)
	Error(w http.ResponseWriter, err error)
}
