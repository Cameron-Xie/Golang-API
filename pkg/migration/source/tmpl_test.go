package source

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTmpl_Open(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "tmpl_test_open")
	if err != nil {
		t.Fatal(err)
	}

	if err := setupTest(tmp); err != nil {
		t.Fatal(err)
	}
	defer teardownTesting(tmp)
	a := assert.New(t)
	m := []struct {
		p   string
		err error
	}{
		{
			p: tmp,
		},
		{
			p:   "random_path",
			err: new(os.PathError),
		},
		{
			p:   ":!",
			err: new(url.Error),
		},
	}

	for _, i := range m {
		d, err := new(tmpl).Open(i.p)
		if err != nil {
			a.IsType(i.err, err)
			continue
		}

		o := &tmpl{path: i.p}

		a.Nil(o.Init(http.Dir(o.path), ""))
		a.Equal(o, d)
	}
}

func TestTmpl_ReadUp(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "tmpl_test_read_up")
	if err != nil {
		t.Fatal(err)
	}
	if err := setupTest(tmp); err != nil {
		t.Fatal(err)
	}
	defer teardownTesting(tmp)

	a := assert.New(t)
	m := []struct {
		dir                string
		data               interface{}
		expectContent      string
		expectedIdentifier string
		err                string
	}{
		{
			dir:                tmp,
			data:               "data",
			expectContent:      "up_input: data",
			expectedIdentifier: "create_schema_and_db_role",
		},
		{
			dir:                path.Join(tmp, "data"),
			data:               "random",
			expectContent:      "up_input: data",
			expectedIdentifier: "create_schema_and_db_role",
			err:                `template: :1:12: executing "" at <.Table>: can't evaluate field Table in type string`,
		},
		{
			dir:                path.Join(tmp, "invalid_template"),
			data:               "data",
			expectContent:      "up_input: data",
			expectedIdentifier: "create_schema_and_db_role",
			err:                `template: :1: unexpected "}" in command`,
		},
	}

	for _, i := range m {
		d, err := (&tmpl{data: i.data}).Open(i.dir)
		if err != nil {
			t.Fatal(err)
		}

		rc, m, err := d.ReadUp(1)
		if err != nil {
			a.Equal(i.err, err.Error())
			a.Nil(rc)
			a.Equal("", m)
			continue
		}

		b, _ := ioutil.ReadAll(rc)
		a.Equal(i.expectContent, string(b))
		a.Equal(i.expectedIdentifier, m)
	}
}

func TestTmpl_read(t *testing.T) {
	a := assert.New(t)
	d := new(tmpl)
	expected := errors.New("something went wrong")
	_, _, err := d.read(nil, "", expected)

	a.Equal(expected, err)
}

func TestNewTmpl(t *testing.T) {
	a := assert.New(t)
	data := "random_string"
	a.Equal(&tmpl{data: data}, NewTmpl(data))
}

func TestTmpl_ReadDown(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "tmpl_test_read_up")
	if err != nil {
		t.Fatal(err)
	}
	if err := setupTest(tmp); err != nil {
		t.Fatal(err)
	}
	defer teardownTesting(tmp)

	a := assert.New(t)
	m := []struct {
		dir                string
		data               interface{}
		expectContent      string
		expectedIdentifier string
		err                string
	}{
		{
			dir:                path.Join(tmp, "invalid_template"),
			data:               "data",
			expectContent:      "down_input: data",
			expectedIdentifier: "create_schema_and_db_role",
			err:                `template: :1: unexpected "}" in command`,
		},
		{
			dir:                path.Join(tmp, "data"),
			data:               "random",
			expectContent:      "down_input: data",
			expectedIdentifier: "create_schema_and_db_role",
			err:                `template: :1:14: executing "" at <.Table>: can't evaluate field Table in type string`,
		},
		{
			dir:                tmp,
			data:               "data",
			expectContent:      "down_input: data",
			expectedIdentifier: "create_schema_and_db_role",
		},
	}

	for _, i := range m {
		s := &tmpl{data: i.data}
		d, err := s.Open(i.dir)
		if err != nil {
			t.Fatal(err)
		}
		rc, m, err := d.ReadDown(1)

		if err != nil {
			a.Equal(i.err, err.Error())
			a.Nil(rc)
			a.Equal("", m)
			continue
		}

		b, _ := ioutil.ReadAll(rc)
		a.Equal(i.expectContent, string(b))
		a.Equal(i.expectedIdentifier, m)
	}
}

func setupTest(tmp string) error {
	m := map[string]string{
		"000001_create_schema_and_db_role.up.sql":                    "up_input: {{.}}",
		"000001_create_schema_and_db_role.down.sql":                  "down_input: {{.}}",
		"invalid_template/000001_create_schema_and_db_role.up.sql":   "up_input: {{}",
		"invalid_template/000001_create_schema_and_db_role.down.sql": "down_input: {{}",
		"data/000001_create_schema_and_db_role.up.sql":               "up_input: {{.Table}}",
		"data/000001_create_schema_and_db_role.down.sql":             "down_input: {{.Table}}",
	}

	if err := createTestingFiles(tmp, m); err != nil {
		return err
	}

	return nil
}

func createTestingFiles(dir string, m map[string]string) error {
	for f, c := range m {
		p := path.Join(dir, f)
		dir, _ := filepath.Split(p)

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		if err := ioutil.WriteFile(p, []byte(c), 06444); err != nil {
			return err
		}
	}

	return nil
}

func teardownTesting(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}
}
