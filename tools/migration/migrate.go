package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Cameron-Xie/Golang-API/pkg/migration/source"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	ms "github.com/golang-migrate/migrate/v4/source"
)

const (
	pgHost             = "POSTGRES_HOST"
	pgPort             = "POSTGRES_PORT"
	admin              = "POSTGRES_USER"
	adminPassword      = "POSTGRES_PASSWORD"
	readWriter         = "POSTGRES_READWRITER"
	readWriterPassword = "POSTGRES_READWRITER_PASSWORD"
	dbName             = "POSTGRES_DB"
)

type Logger struct {
	V bool
}

func (l *Logger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *Logger) Verbose() bool {
	return l.V
}

type Cred struct {
	Username string
	Password string
}

func main() {
	fmt.Println("Start migration...")

	ms.Register(source.TmplSchema, source.NewTmpl(Cred{
		Username: os.Getenv(readWriter),
		Password: os.Getenv(readWriterPassword),
	}))

	dir, err := os.Getwd()
	handleError(err)

	p := filepath.Join(dir, "database/migrations")
	m, err := migrate.New(
		fmt.Sprintf("%v://%v", source.TmplSchema, p),
		fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v?sslmode=require",
			os.Getenv(admin),
			url.QueryEscape(os.Getenv(adminPassword)),
			os.Getenv(pgHost),
			os.Getenv(pgPort),
			os.Getenv(dbName),
		),
	)
	handleError(err)

	m.Log = &Logger{V: true}
	handleError(m.Up())
	fmt.Println("Migration Success")
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
