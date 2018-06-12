package handler

import (
	"net/http"
	"encoding/json"
)

type JsonResponse struct {
	Writer http.ResponseWriter
}

func NewJsonResponse(w http.ResponseWriter) *JsonResponse {
	return (&JsonResponse{Writer: w}).setHeader()
}

func (res *JsonResponse) setHeader() *JsonResponse {
	res.Writer.Header().Set("Content-Type", "application/json")
	return res
}

func (res *JsonResponse) Response(payload interface{}, statusCode int, isSuccess bool) error {

	payloadKey := "data"

	if isSuccess != true {
		payloadKey = "error"
	}

	jsonBytes, err := json.Marshal(map[string]interface{}{payloadKey: payload})

	if err != nil {
		return err
	}

	res.Writer.Write(jsonBytes)

	return nil
}
