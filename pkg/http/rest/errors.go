package rest

import (
	"encoding/json"
	"net/http"
)

const (
	notFoundErr string = "item not found"
)

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type HTTPError struct {
	Title         string         `json:"title"`
	InvalidParams []InvalidParam `json:"invalid_params,omitempty"`
	StatusCode    int            `json:"-"`
}

func (e *HTTPError) Error() string {
	b, _ := json.Marshal(e)

	return string(b)
}

func NotFoundError() *HTTPError {
	return &HTTPError{
		Title:      notFoundErr,
		StatusCode: http.StatusNotFound,
	}
}
