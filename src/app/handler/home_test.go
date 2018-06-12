package handler

import (
	"testing"
	"net/http"
	"github.com/gavv/httpexpect"
)

func TestHomeHandler_Welcome(t *testing.T) {
	home := &HomeHandler{}
	server := createTestServer("/", home.Welcome)
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("data").ValueEqual("data", "Hello World")
}
