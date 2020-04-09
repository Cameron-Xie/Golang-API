package source

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
)

const TmplSchema string = "tmpl"

type tmpl struct {
	httpfs.PartialDriver
	path string
	data interface{}
}

func NewTmpl(data interface{}) source.Driver {
	return &tmpl{data: data}
}

func (s *tmpl) Open(i string) (source.Driver, error) {
	u, err := url.Parse(i)
	if err != nil {
		return nil, err
	}

	s.path = filepath.Join(u.Host, u.Path)
	if err := s.Init(http.Dir(s.path), ""); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *tmpl) ReadUp(v uint) (io.ReadCloser, string, error) {
	return s.read(s.PartialDriver.ReadUp(v))
}

func (s *tmpl) ReadDown(v uint) (io.ReadCloser, string, error) {
	return s.read(s.PartialDriver.ReadDown(v))
}

func (s *tmpl) read(rc io.ReadCloser, i string, err error) (io.ReadCloser, string, error) {
	if err != nil {
		return rc, i, err
	}

	defer rc.Close()
	b, _ := ioutil.ReadAll(rc)

	t, err := template.New("").Parse(string(b))
	if err != nil {
		return nil, "", err
	}

	var rs strings.Builder
	err = t.Execute(&rs, s.data)
	if err != nil {
		return nil, "", err
	}

	return ioutil.NopCloser(strings.NewReader(rs.String())), i, nil
}
