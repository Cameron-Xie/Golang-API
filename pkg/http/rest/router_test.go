package rest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	a := assert.New(t)
	r := NewRouter(
		map[string]http.Handler{
			"/hello": new(helloWorldMockHandler),
		},
	)

	mws := r.(*chi.Mux).Middlewares()
	ts := httptest.NewServer(r)
	m := []struct {
		m       string
		p       string
		resBody string
	}{
		{m: "GET", p: "/hello", resBody: "Hello World!"},
	}

	for _, s := range m {
		resp := getHTTPRespStr(t, ts, s.m, s.p, nil)
		a.Equal(s.resBody, resp)
	}

	// checking middleware
	fs := []string{
		"github.com/go-chi/chi/middleware.RequestID",
		"github.com/go-chi/chi/middleware.Recoverer",
	}

	for i := 0; i < len(fs); i++ {
		a.Equal(fs[i], getFuncName(mws[i]))
	}
}

type helloWorldMockHandler struct {
}

func (h *helloWorldMockHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

func getHTTPRespStr(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) string {
	req, err := http.NewRequest(method, fmt.Sprintf("%v%v", ts.URL, path), body)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	return string(resBody)
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
