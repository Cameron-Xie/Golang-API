package handler

import (
	"net/http"
)

type HomeHandler struct{}

type Payload struct {
	Message string
}

func (handler *HomeHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	NewJsonResponse(w).Response("Hello World", http.StatusOK, true)
}
