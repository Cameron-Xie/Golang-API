package main

import (
	"log"
	"net/http"
	"os"
	"time"

	restI "github.com/Cameron-Xie/Golang-API/internal/http/rest"
	"github.com/Cameron-Xie/Golang-API/internal/storage/postgres"
	"github.com/Cameron-Xie/Golang-API/pkg/http/rest"
	"github.com/Cameron-Xie/Golang-API/pkg/logger"
	"github.com/Cameron-Xie/Golang-API/pkg/services/readtask"
	"github.com/Cameron-Xie/Golang-API/pkg/services/storetask"
	"github.com/Cameron-Xie/Golang-API/pkg/services/updatetask"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	pgHost             string = "POSTGRES_HOST"
	pgPort             string = "POSTGRES_PORT"
	readWriter         string = "POSTGRES_READWRITER"
	readWriterPassword string = "POSTGRES_READWRITER_PASSWORD" // nolint: gosec
	dbName             string = "POSTGRES_DB"
	sslMode            string = "require" // "verify-full"
	maxDBConnections   int    = 10
	serverAddr         string = ":8080"
	serverReadTimeout         = time.Second
	serverWriteTimeout        = 5 * time.Second
	pageLimit          int    = 100
	version            int    = 1
)

func main() {
	storage, err := postgres.New(&postgres.DBConn{
		Host:     os.Getenv(pgHost),
		Port:     os.Getenv(pgPort),
		Database: os.Getenv(dbName),
		Username: os.Getenv(readWriter),
		Password: os.Getenv(readWriterPassword),
		SSLMode:  sslMode,
		MaxConn:  maxDBConnections,
		ConnLife: time.Minute,
	}).Open()
	errHandler(err)

	r := rest.NewRouter(
		version,
		map[string]http.Handler{
			"/tasks": rest.NewEndpoint(
				storetask.New(storetask.NewValidator(), storage),
				updatetask.New(updatetask.NewValidator(), storage),
				readtask.New(storage),
				restI.NewJSONRender(),
				pageLimit,
			),
		},
		rest.NewLogFormatter(logger.NewJSONLogger()),
		rest.SetCORS([]string{"*"}),
	)

	log.Fatal(setServer(serverAddr, r).ListenAndServe())
}

func setServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
	}
}

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
