package storetask

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Store(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		input       io.ReadCloser
		validateErr error
		storageErr  error
		expected    *Task
		expectedErr string
	}{
		{
			input:       ioutil.NopCloser(bytes.NewBuffer([]byte(`invalid_json`))),
			expectedErr: "invalid character 'i' looking for beginning of value",
		},
		{
			input:       ioutil.NopCloser(bytes.NewBuffer([]byte(`{"Name":"task_a"}`))),
			validateErr: errors.New("something went wrong"),
			expectedErr: "something went wrong",
		},
		{
			input:       ioutil.NopCloser(bytes.NewBuffer([]byte(`{"Name":"task_a"}`))),
			storageErr:  errors.New("something went wrong"),
			expectedErr: "something went wrong",
		},
		{
			input:    ioutil.NopCloser(bytes.NewBuffer([]byte(`{"Name":"task_a"}`))),
			expected: &Task{Name: "task_a"},
		},
	}

	for _, i := range m {
		taskType := reflect.TypeOf(new(Task)).String()
		v := new(validatorMock)
		s := new(storageMock)
		v.On("Validate", mock.AnythingOfType(taskType)).Return(i.validateErr)
		s.On("StoreTask", mock.AnythingOfType(taskType)).Return(i.storageErr)

		svc := service{
			storage:   s,
			validator: v,
		}
		res, err := svc.Store(i.input)

		if i.expectedErr == "" {
			task := res.(*Task)
			a.True(isUUID(task.ID.String(), 4))
			a.Equal(i.expected.Name, task.Name)
			a.Equal(i.expected.Description, task.Description)
			a.Nil(err)

			continue
		}

		a.Nil(i.expected)
		a.NotNil(err)
		a.Equal(i.expectedErr, err.Error())
	}
}

func TestService_Delete(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		storageErr error
	}{
		{
			storageErr: errors.New("delete failed"),
		},
	}

	for _, i := range m {
		idType := reflect.TypeOf(uuid.UUID{}).String()
		storage := new(storageMock)
		s := &service{storage: storage}

		storage.On("DeleteTask", mock.AnythingOfType(idType)).Return(i.storageErr)

		a.Equal(i.storageErr, s.Delete(uuid.New()))
	}
}

func TestNew(t *testing.T) {
	a := assert.New(t)
	v := new(validatorMock)
	s := new(storageMock)

	a.Equal(&service{
		storage:   s,
		validator: v,
	}, New(v, s))
}

type storageMock struct {
	mock.Mock
}

func (s *storageMock) StoreTask(t *Task) error {
	args := s.Called(t)

	return args.Error(0)
}

func (s *storageMock) DeleteTask(id uuid.UUID) error {
	args := s.Called(id)

	return args.Error(0)
}

type validatorMock struct {
	mock.Mock
}

func (v *validatorMock) Validate(t *Task) error {
	args := v.Called(t)

	return args.Error(0)
}

func isUUID(id string, v int) bool {
	i, err := uuid.Parse(id)

	if err != nil || i.Version() != uuid.Version(v) {
		return false
	}

	return true
}
