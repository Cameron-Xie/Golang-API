package readtask

import (
	"github.com/Cameron-Xie/Golang-API/pkg/http/rest"
	"github.com/google/uuid"
)

type Storage interface {
	FetchTask(uuid.UUID) (*Task, error)
	FetchTasks(offset, limit int) (*rest.Collection, error)
}

type service struct {
	storage Storage
}

func (s *service) List(offset, limit int) (*rest.Collection, error) {
	return s.storage.FetchTasks(offset, limit)
}

func (s *service) Read(id uuid.UUID) (interface{}, error) {
	return s.storage.FetchTask(id)
}

func New(s Storage) rest.ReadService {
	return &service{storage: s}
}
