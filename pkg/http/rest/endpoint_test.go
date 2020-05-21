package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// nolint funlen
func TestNewEndpoint(t *testing.T) {
	a := assert.New(t)
	id := uuid.New()
	m := []endpointTestRequest{
		// GET / 200
		{
			method: http.MethodGet,
			resp: fmt.Sprintf(
				`{"meta":{"total":2},"links":[{"href":"/?limit=10\u0026page=1","rel":"self"},{"href":"/?limit=10\u0026page=1","rel":"first"},{"href":"/?limit=10\u0026page=1","rel":"last"}],"items":[{"id":"%v","name":"read_svc_a"},{"id":"%v","name":"read_svc_b"}]}`,
				id.String(),
				id.String(),
			),
			code: 200,
		},
		// GET /?page=1&limit=1 200
		{
			method: http.MethodGet,
			path:   "?page=1&limit=1",
			resp: fmt.Sprintf(
				`{"meta":{"total":2},"links":[{"href":"/?limit=1\u0026page=1","rel":"self"},{"href":"/?limit=1\u0026page=1","rel":"first"},{"href":"/?limit=1\u0026page=2","rel":"next"},{"href":"/?limit=1\u0026page=2","rel":"last"}],"items":[{"id":"%v","name":"read_svc_a"}]}`,
				id.String(),
			),
			code: 200,
		},
		// GET /?page=2&limit=1 200
		{
			method: http.MethodGet,
			path:   "?page=2&limit=1",
			resp: fmt.Sprintf(
				`{"meta":{"total":2},"links":[{"href":"/?limit=1\u0026page=2","rel":"self"},{"href":"/?limit=1\u0026page=1","rel":"first"},{"href":"/?limit=1\u0026page=1","rel":"prev"},{"href":"/?limit=1\u0026page=2","rel":"last"}],"items":[{"id":"%v","name":"read_svc_b"}]}`,
				id.String(),
			),
			code: 200,
		},
		// GET / error
		{
			method: http.MethodGet,
			path:   "?page=10",
			resp:   `out of range`,
			code:   400,
		},
		// POST / 201
		{
			method:    http.MethodPost,
			body:      strings.NewReader(`{}`),
			isContain: true,
			resp:      `"id"`,
			code:      201,
		},
		// POST / 400
		{
			method:    http.MethodPost,
			body:      strings.NewReader(`random`),
			isContain: true,
			resp:      `invalid character`,
			code:      400,
		},
		// GET /{ID} 200
		{
			method: http.MethodGet,
			path:   id.String(),
			resp: fmt.Sprintf(
				`{"id":"%v","name":"read_svc_a"}`,
				id.String(),
			),
			code: 200,
		},
		// GET /{ID} error
		{
			method: http.MethodGet,
			path:   uuid.New().String(),
			resp:   `not found`,
			code:   400,
		},
		// PATCH /{ID} 204
		{
			method: http.MethodPatch,
			path:   id.String(),
			code:   204,
		},
		// PATCH /{invalid_uuid} error
		{
			method: http.MethodPatch,
			path:   "invalid_uuid",
			resp:   `{"title":"item not found"}`,
			code:   400,
		},
		// PATCH /{ID} error
		{
			method: http.MethodPatch,
			path:   uuid.New().String(),
			resp:   `not found`,
			code:   400,
		},
		// DELETE /{ID} 204
		{
			method: http.MethodDelete,
			path:   id.String(),
			code:   204,
		},
		// DELETE /{ID} error
		{
			method: http.MethodDelete,
			path:   uuid.New().String(),
			resp:   `not found`,
			code:   400,
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(m))

	for _, i := range m {
		go func(r endpointTestRequest) {
			defer wg.Done()
			ts := setTestServer(id)

			req, _ := http.NewRequest(
				r.method,
				fmt.Sprintf("%v/%v", ts.URL, r.path),
				r.body,
			)
			req.Header.Set("Content-Type", "application/json")

			resp, err := new(http.Client).Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer func() { _ = resp.Body.Close() }()
			b, _ := ioutil.ReadAll(resp.Body)

			a.Equal(r.code, resp.StatusCode)

			if r.isContain {
				a.Contains(string(b), r.resp)
			} else {
				a.Equal(r.resp, string(b))
			}
		}(i)
	}

	wg.Wait()
}

func setTestServer(id uuid.UUID) *httptest.Server {
	storeSvc := &storeSvcMock{
		items: []endpointTestItem{
			{
				ID: id,
			},
		},
	}

	updateSvc := &updateSvcMock{
		items: []endpointTestItem{
			{
				ID: id,
			},
		},
	}

	readSvc := &readSvcMock{
		items: []endpointTestItem{
			{
				ID:   id,
				Name: "read_svc_a",
			},
			{
				ID:   id,
				Name: "read_svc_b",
			},
		},
	}

	return httptest.NewServer(
		NewEndpoint(
			storeSvc,
			updateSvc,
			readSvc,
			new(renderMock),
			10,
		),
	)
}

type endpointTestItem struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type storeSvcMock struct {
	items []endpointTestItem
}

func (s *storeSvcMock) Store(r io.ReadCloser) (interface{}, error) {
	i := new(endpointTestItem)
	if err := json.NewDecoder(r).Decode(i); err != nil {
		return nil, err
	}

	return i, nil
}

func (s *storeSvcMock) Update(_ io.ReadCloser, i uuid.UUID) error {
	for _, item := range s.items {
		if item.ID == i {
			return nil
		}
	}

	return errors.New("not found")
}

func (s *storeSvcMock) Delete(id uuid.UUID) error {
	for _, item := range s.items {
		if item.ID == id {
			return nil
		}
	}

	return errors.New("not found")
}

type updateSvcMock struct {
	items []endpointTestItem
}

func (s *updateSvcMock) Update(_ io.ReadCloser, i uuid.UUID) error {
	for _, item := range s.items {
		if item.ID == i {
			return nil
		}
	}

	return errors.New("not found")
}

type readSvcMock struct {
	items []endpointTestItem
}

func (s *readSvcMock) List(offset, limit int) (coll *Collection, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("out of range")
		}
	}()

	return &Collection{
		Total: len(s.items),
		Items: s.items[offset:min(offset+limit, len(s.items))],
	}, nil
}

func (s *readSvcMock) Read(id uuid.UUID) (interface{}, error) {
	for _, item := range s.items {
		if item.ID == id {
			return item, nil
		}
	}

	return nil, errors.New("not found")
}

type renderMock struct {
}

func (r *renderMock) Success(w http.ResponseWriter, c interface{}, code int) {
	w.WriteHeader(code)
	b, _ := json.Marshal(c)

	_, _ = w.Write(b)
}

func (r *renderMock) Error(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(err.Error()))
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}

type endpointTestRequest struct {
	method    string
	path      string
	body      io.Reader
	isContain bool
	resp      string
	code      int
}
