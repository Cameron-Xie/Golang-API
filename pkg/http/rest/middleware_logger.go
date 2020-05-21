package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Cameron-Xie/Golang-API/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type logEntryCreator struct {
	Logger logger.Logger
}

type logEntry struct {
	Logger logger.Logger
}

func (e *logEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	e.Logger = e.Logger.WithFields(logrus.Fields{
		"respStatus":      status,
		"respBytesLength": bytes,
		"respElapsedMS":   float64(elapsed.Nanoseconds()) / float64(time.Millisecond),
	})

	e.Logger.Infoln("end request")
}

func (e *logEntry) Panic(v interface{}, stack []byte) {
	e.Logger = e.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func (c *logEntryCreator) NewLogEntry(r *http.Request) middleware.LogEntry {
	e := &logEntry{Logger: c.Logger}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	e.Logger = e.Logger.WithFields(
		logrus.Fields{
			"timestamp":  time.Now().Format(time.RFC3339),
			"httpScheme": scheme,
			"httpProto":  r.Proto,
			"httpMethod": r.Method,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
			"uri":        fmt.Sprintf("%s://%s%v", scheme, r.Host, r.RequestURI),
		},
	)

	e.Logger.Infoln("start request")

	return e
}

func NewLogFormatter(l logger.Logger) middleware.LogFormatter {
	return &logEntryCreator{Logger: l}
}
