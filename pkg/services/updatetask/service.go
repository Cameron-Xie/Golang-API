package updatetask

import (
	"encoding/json"
	"io"

	"github.com/Cameron-Xie/Golang-API/pkg/http/rest"
	"github.com/google/uuid"
)

type Validator interface {
	Validate(map[string]interface{}) (map[string]interface{}, error)
}

type Storage interface {
	UpdateTask(map[string]interface{}, uuid.UUID) error
}

type service struct {
	validator Validator
	storage   Storage
}

func (s *service) Update(r io.ReadCloser, id uuid.UUID) error {
	defer func() { _ = r.Close() }()

	i := make(map[string]interface{})
	if err := json.NewDecoder(r).Decode(&i); err != nil {
		return err
	}

	res, err := s.validator.Validate(i)
	if err != nil {
		return err
	}

	if err := s.storage.UpdateTask(res, id); err != nil {
		return err
	}

	return nil
}

func New(v Validator, s Storage) rest.UpdateService {
	return &service{
		validator: v,
		storage:   s,
	}
}
