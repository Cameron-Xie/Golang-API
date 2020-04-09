package updatetask

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

func TestService_Update(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		input       io.ReadCloser
		validateErr error
		storageErr  error
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
			input: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"Name":"task_a"}`))),
		},
	}

	for _, i := range m {
		inputType := reflect.TypeOf(make(map[string]interface{})).String()
		idType := reflect.TypeOf(uuid.UUID{}).String()
		v := new(validatorMock)
		s := new(storageMock)
		v.On("Validate", mock.AnythingOfType(inputType)).Return(i.validateErr)
		s.On(
			"UpdateTask",
			mock.AnythingOfType(inputType),
			mock.AnythingOfType(idType),
		).Return(i.storageErr)

		svc := service{
			storage:   s,
			validator: v,
		}

		err := svc.Update(i.input, uuid.New())
		if i.expectedErr == "" {
			a.Nil(err)
		} else {
			a.NotNil(err)
			a.Equal(i.expectedErr, err.Error())
		}
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

func (s *storageMock) UpdateTask(i map[string]interface{}, id uuid.UUID) error {
	args := s.Called(i, id)

	return args.Error(0)
}

type validatorMock struct {
	mock.Mock
}

func (v *validatorMock) Validate(i map[string]interface{}) (map[string]interface{}, error) {
	args := v.Called(i)

	return i, args.Error(0)
}
