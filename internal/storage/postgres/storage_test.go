package postgres

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/stretchr/testify/assert"
)

const (
	pgHost             = "POSTGRES_HOST"
	pgPort             = "POSTGRES_PORT"
	readWriter         = "POSTGRES_READWRITER"
	readWriterPassword = "POSTGRES_READWRITER_PASSWORD" // nolint: gosec
	dbName             = "POSTGRES_DB"
)

func TestNewStorage(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		conn *DBConn
		err  string
	}{
		{
			conn: getValidTestConn(),
		},
		{
			conn: &DBConn{},
			err:  "dial tcp: lookup port=: no such host",
		},
	}

	for _, i := range m {
		_, err := New(i.conn).Open()

		if i.err != "" {
			a.NotNil(err)
			a.Equal(i.err, err.Error())
		} else {
			a.Nil(err)
		}
	}
}

func getValidTestConn() *DBConn {
	return &DBConn{
		Host:     os.Getenv(pgHost),
		Port:     os.Getenv(pgPort),
		Database: os.Getenv(dbName),
		Username: os.Getenv(readWriter),
		Password: os.Getenv(readWriterPassword),
		MaxConn:  10,
		ConnLife: time.Minute,
	}
}

func contains(l []string, s string) bool {
	sort.Strings(l)
	i := sort.SearchStrings(l, s)
	return i < len(l) && l[i] == s
}

// compare Exposed Fields in struct
func isEqual(expected, target interface{}, excepts []string) error {
	ve := reflect.ValueOf(expected)
	te := reflect.TypeOf(expected)
	vt := reflect.ValueOf(target)

	if ve.Type() != vt.Type() {
		return fmt.Errorf("expected type %v is not equal to target type %v", ve.Type(), vt.Type())
	}

	if ve.Kind() == reflect.Ptr {
		ve = ve.Elem()
		te = te.Elem()
		vt = vt.Elem()
	}

	if ve.Kind() != reflect.Struct {
		return errors.New("not struct")
	}

	for i := 0; i < ve.NumField(); i++ {
		n := te.Field(i).Name
		if contains(excepts, n) {
			continue
		}

		fe := ve.FieldByName(n)
		ft := vt.FieldByName(n)

		if !fe.CanInterface() {
			continue
		}

		if fe.Interface() != ft.Interface() {
			return fmt.Errorf("field %v: expected %v, but got %v", n, fe.Interface(), ft.Interface())
		}
	}

	return nil
}
