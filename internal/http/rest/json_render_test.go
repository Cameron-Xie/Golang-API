package rest

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Cameron-Xie/Golang-API/internal/storage/postgres"
	"github.com/Cameron-Xie/Golang-API/pkg/http/rest"
	"github.com/stretchr/testify/assert"
)

func TestJsonRender_Success(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		content      interface{}
		code         int
		expected     string
		expectedCode int
	}{
		{
			content:      make(chan int),
			code:         200,
			expected:     `{"title":"json: unsupported type: chan int"}`,
			expectedCode: 502,
		},
		{
			content:      "random_string",
			code:         200,
			expected:     `{"title":"invalid json"}`,
			expectedCode: 502,
		},
		{
			content:      `{"name":"user_name"}`,
			code:         200,
			expected:     `{"name":"user_name"}`,
			expectedCode: 200,
		},
		{
			content:      []byte(`{"name":"user_name"}`),
			code:         200,
			expected:     `{"name":"user_name"}`,
			expectedCode: 200,
		},
		{
			content: &testRespStruct{
				Name: "username",
			},
			code:         200,
			expected:     `{"Name":"username"}`,
			expectedCode: 200,
		},
	}

	for _, i := range m {
		resp := newRespWriteMock()
		w := NewJSONRender()
		w.Success(resp, i.content, i.code)

		a.Equal(i.expected, string(resp.body))
		a.Equal(i.expectedCode, resp.code)
		a.Equal("application/json", resp.header.Get("Content-Type"))
	}
}

func TestJsonRender_Error(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		err          error
		expected     string
		expectedCode int
	}{
		{
			err:          errors.New("something went wrong"),
			expected:     `{"title":"something went wrong"}`,
			expectedCode: 502,
		},
		{
			err: &rest.HTTPError{
				Title: "title",
				InvalidParams: []rest.InvalidParam{
					{
						Name:   "ID",
						Reason: "invalid",
					},
				},
				StatusCode: 400,
			},
			expected:     `{"title":"title","invalid_params":[{"name":"ID","reason":"invalid"}]}`,
			expectedCode: 400,
		},
		{
			err: &postgres.NotFoundError{
				Table: "table_name",
				Value: "id_value",
			},
			expected:     `{"title":"id_value is not found in table_name"}`,
			expectedCode: 404,
		},
	}

	for _, i := range m {
		resp := newRespWriteMock()
		w := NewJSONRender()
		w.Error(resp, i.err)

		a.Equal(i.expected, string(resp.body))
		a.Equal(i.expectedCode, resp.code)
		a.Equal("application/json", resp.header.Get("Content-Type"))
	}
}

type testRespStruct struct {
	Name string
}

type respWriterMock struct {
	code   int
	body   []byte
	header http.Header
}

func (w *respWriterMock) Header() http.Header {
	return w.header
}

func (w *respWriterMock) Write(p []byte) (int, error) {
	w.body = p
	return len(p), nil
}

func (w *respWriterMock) WriteHeader(statusCode int) {
	w.code = statusCode
}

func newRespWriteMock() *respWriterMock {
	return &respWriterMock{
		code:   0,
		body:   nil,
		header: make(http.Header),
	}
}
