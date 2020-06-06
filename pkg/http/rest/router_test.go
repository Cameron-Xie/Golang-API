package rest

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	type routeTest struct {
		m         string
		p         string
		reqHeader http.Header
		resBody   string
		resHeader http.Header
	}

	a := assert.New(t)
	r := NewRouter(
		1,
		map[string]http.Handler{
			"/hello": new(helloWorldMockHandler),
		},
		new(loggerEntryCreatorMock),
		SetCORS([]string{"https://example.com"}),
		middleware.NewCompressor(5),
	)

	mws := r.(*chi.Mux).Middlewares()
	m := []routeTest{
		{
			m: "OPTIONS",
			p: "/v1/hello",
			reqHeader: map[string][]string{
				"Origin":                        {"https://example.com"},
				"Access-Control-Request-Method": {"GET"},
			},
			resBody: "",
			resHeader: map[string][]string{
				"Access-Control-Allow-Credentials": {"true"},
				"Access-Control-Allow-Origin":      {"https://example.com"},
				"Access-Control-Allow-Methods":     {"GET"},
			},
		},
		{m: "GET", p: "/v1/hello", resBody: "Hello World!"},
	}

	var wg sync.WaitGroup
	wg.Add(len(m))

	for _, s := range m {
		go func(i routeTest) {
			defer wg.Done()
			ts := httptest.NewServer(r)
			res := testHTTPRequest(ts, i.m, i.p, nil, i.reqHeader)
			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			res.Body.Close()

			for k := range i.resHeader {
				a.Equal(i.resHeader.Get(k), res.Header.Get(k))
			}

			a.Equal(i.resBody, string(resBody))
		}(s)
	}

	// checking middleware
	fs := []string{
		"github.com/go-chi/chi/middleware.RequestID",
		"github.com/go-chi/chi/middleware.RequestLogger.func1",
		"github.com/go-chi/chi/middleware.Recoverer",
		"github.com/go-chi/chi/middleware.Compress",
	}

	for i := 0; i < len(fs); i++ {
		a.Equal(fs[i], getFuncName(mws[i]))
	}

	wg.Wait()
}

type loggerEntryCreatorMock struct {
}

type loggerEntryMock struct {
	middleware.LogEntry
}

func (e *loggerEntryMock) Write(_, _ int, _ http.Header, _ time.Duration, _ interface{}) {
}

func (c *loggerEntryCreatorMock) NewLogEntry(_ *http.Request) middleware.LogEntry {
	return new(loggerEntryMock)
}

type helloWorldMockHandler struct {
}

func (h *helloWorldMockHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

func testHTTPRequest(
	ts *httptest.Server,
	method, path string,
	body io.Reader,
	headers http.Header) *http.Response {
	req, err := http.NewRequest(method, fmt.Sprintf("%v%v", ts.URL, path), body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header = headers

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
