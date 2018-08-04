package handler

import (
	"net/http"
	"encoding/json"
)

func JsonResponse(w http.ResponseWriter) func(*Response) {
	return func(response *Response) {
		payloadKey := "data"

		if response.isSuccess == false {
			payloadKey = "error"
		}

		payload, _ := getJsonPayload(response.payload, payloadKey)

		writeResponse(setJsonHeader(w))(payload, response.code)
	}
}

func setJsonHeader(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	return w
}

func getJsonPayload(payload interface{}, payloadKey string) ([]byte, error) {
	return json.Marshal(map[string]interface{}{payloadKey: payload})
}
