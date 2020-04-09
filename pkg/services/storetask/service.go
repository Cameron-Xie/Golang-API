package storetask

import (
	"encoding/json"
	"io"

	"github.com/Cameron-Xie/Golang-API/pkg/http/rest"
	"github.com/google/uuid"
)

type Validator interface {
	Validate(*Task) error
}

type Storage interface {
	StoreTask(*Task) error
	DeleteTask(uuid.UUID) error
}

type service struct {
	validator Validator
	storage   Storage
}

func New(v Validator, s Storage) rest.StoreService {
	return &service{
		storage:   s,
		validator: v,
	}
}

func (s *service) Store(r io.ReadCloser) (interface{}, error) {
	defer func() { _ = r.Close() }()
	i := new(Task)

	if err := json.NewDecoder(r).Decode(i); err != nil {
		return nil, err
	}

	if err := s.validator.Validate(i); err != nil {
		return nil, err
	}

	i.ID = uuid.New()
	if err := s.storage.StoreTask(i); err != nil {
		return nil, err
	}

	return i, nil
}

func (s *service) Delete(id uuid.UUID) error {
	return s.storage.DeleteTask(id)
}
