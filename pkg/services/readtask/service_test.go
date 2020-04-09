package readtask

import (
	"errors"
	"testing"

	"github.com/Cameron-Xie/Golang-API/pkg/http/rest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_List(t *testing.T) {
	a := assert.New(t)
	id := uuid.New()
	m := []struct {
		offset   int
		limit    int
		expected *rest.Collection
		err      error
	}{
		{
			offset: 0,
			limit:  10,
			expected: &rest.Collection{
				Total: 1,
				Items: []Task{
					{
						ID:   id,
						Name: "task_a",
					},
				},
			},
		},
		{
			offset:   10,
			limit:    10,
			expected: nil,
			err:      errors.New("out of range"),
		},
	}

	for _, i := range m {
		s := &service{storage: newStorageMock(id)}
		coll, err := s.List(i.offset, i.limit)

		a.Equal(i.err, err)
		a.Equal(i.expected, coll)
	}
}

func TestService_Read(t *testing.T) {
	a := assert.New(t)
	id := uuid.New()
	m := []struct {
		id       uuid.UUID
		expected *Task
		err      error
	}{
		{
			id: id,
			expected: &Task{
				ID:   id,
				Name: "task_a",
			},
		},
		{
			id:  uuid.New(),
			err: errors.New("not found"),
		},
	}

	for _, i := range m {
		s := &service{storage: newStorageMock(id)}
		res, err := s.Read(i.id)

		a.Equal(i.err, err)
		if err == nil {
			a.Equal(i.expected, res.(*Task))
		} else {
			a.Nil(res)
		}
	}
}

func TestNew(t *testing.T) {
	a := assert.New(t)
	s := newStorageMock(uuid.New())
	a.Equal(&service{storage: s}, New(s))
}

type storageMock struct {
	tasks []Task
}

func (s *storageMock) FetchTask(id uuid.UUID) (*Task, error) {
	for _, i := range s.tasks {
		if i.ID == id {
			return &i, nil
		}
	}

	return nil, errors.New("not found")
}

func (s *storageMock) FetchTasks(offset, limit int) (coll *rest.Collection, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("out of range")
		}
	}()

	return &rest.Collection{
		Total: len(s.tasks),
		Items: s.tasks[offset:min(offset+limit, len(s.tasks))],
	}, nil
}

func newStorageMock(id uuid.UUID) Storage {
	return &storageMock{
		tasks: []Task{
			{
				ID:   id,
				Name: "task_a",
			},
		},
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}
