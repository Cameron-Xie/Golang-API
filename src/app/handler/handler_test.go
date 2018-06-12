package handler

import (
	"net/http"
	"net/http/httptest"
)

func createTestServer(path string, handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(path, handler)
	server := httptest.NewServer(mux)

	return server
}

