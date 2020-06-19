package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Cameron-Xie/Golang-API/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggerEntryCreator_NewLogEntry(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		r          *http.Request
		scheme     string
		proto      string
		method     string
		remoteAddr string
		uri        string
		ignore     bool
	}{
		{
			r: &http.Request{
				Method:     "GET",
				Proto:      "HTTP/1.1",
				Close:      false,
				Host:       "8.8.8.8:8080",
				RemoteAddr: "8.8.8.8",
				RequestURI: "/health",
				URL:        &url.URL{Path: "/health"},
				TLS:        new(tls.ConnectionState),
			},
			scheme:     "https",
			proto:      "HTTP/1.1",
			method:     "GET",
			remoteAddr: "8.8.8.8",
			uri:        "https://8.8.8.8:8080/health",
		},
		{
			r: &http.Request{
				Method:     "GET",
				Proto:      "HTTP/1.1",
				Close:      false,
				Host:       "8.8.8.8:8080",
				RemoteAddr: "8.8.8.8",
				RequestURI: "/health?search=param",
				URL:        &url.URL{Path: "/health"},
				TLS:        nil,
			},
			scheme:     "http",
			proto:      "HTTP/1.1",
			method:     "GET",
			remoteAddr: "8.8.8.8",
			uri:        "http://8.8.8.8:8080/health?search=param",
		},
		{
			r: &http.Request{
				Method:     "GET",
				Proto:      "HTTP/1.1",
				Close:      false,
				Host:       "8.8.8.8:8080",
				RemoteAddr: "8.8.8.8",
				RequestURI: "/ignore",
				URL:        &url.URL{Path: "/ignore"},
				TLS:        nil,
			},
			ignore: true,
		},
	}

	for _, i := range m {
		l := logrus.New()
		l.Out = &bytes.Buffer{}
		l.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: true,
		}
		c := &logEntryCreator{Logger: l, exclude: []string{"/ignore"}}
		c.NewLogEntry(i.r)

		if i.ignore {
			a.Empty(l.Out)
			continue
		}

		var o map[string]string
		_ = json.Unmarshal([]byte(fmt.Sprint(l.Out)), &o)

		_, err := time.Parse(time.RFC3339, o["timestamp"])
		a.Nil(err)
		a.Equal(i.scheme, o["httpScheme"])
		a.Equal(i.proto, o["httpProto"])
		a.Equal(i.method, o["httpMethod"])
		a.Equal(i.remoteAddr, o["remoteAddr"])
		a.Equal(i.uri, o["uri"])
		a.Equal("info", o["level"])
		a.Equal("start request", o["msg"])
	}
}

func TestLoggerEntry_Write(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		s       int
		b       int
		e       time.Duration
		disable bool
	}{
		{
			s: 200,
			b: 1000000.0,
			e: time.Duration(100),
		},
		{
			disable: true,
		},
	}

	for _, i := range m {
		l := logrus.New()
		l.Out = &bytes.Buffer{}
		l.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: true,
		}
		e := &logEntry{Logger: l}
		e.disable = i.disable
		e.Write(i.s, i.b, http.Header{}, i.e, nil)

		if i.disable {
			a.Empty(l.Out)
			continue
		}

		o := &struct {
			Status  int     `json:"respStatus"`
			Size    int     `json:"respBytesLength"`
			Elapsed float64 `json:"respElapsedMS"`
		}{}
		_ = json.Unmarshal([]byte(fmt.Sprint(l.Out)), o)

		a.Equal(i.s, o.Status)
		a.Equal(i.b, o.Size)
		a.Equal(float64(i.e.Nanoseconds())/1000000.0, o.Elapsed)
	}
}

func TestLoggerEntry_Panic(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		v interface{}
		s []byte
	}{
		{
			v: "panic info",
			s: []byte("stacks"),
		},
		{
			v: nil,
			s: []byte("stacks"),
		},
	}

	for _, i := range m {
		l := logrus.New()
		l.Out = &bytes.Buffer{}
		l.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: true,
		}
		e := &logEntry{Logger: l}
		e.Panic(i.v, i.s)
		e.Logger.Print()

		o := make(map[string]string)
		_ = json.Unmarshal([]byte(fmt.Sprint(l.Out)), &o)

		a.Equal(fmt.Sprintf("%+v", i.v), o["panic"])
		a.Equal(string(i.s), o["stack"])
	}
}

func TestNewStdoutJsonLogger(t *testing.T) {
	a := assert.New(t)
	l := logger.NewJSONLogger()
	c := NewLogFormatter(l, []string{"/"}).(*logEntryCreator)

	a.Equal(l, c.Logger)
	a.Equal([]string{"/"}, c.exclude)
}
