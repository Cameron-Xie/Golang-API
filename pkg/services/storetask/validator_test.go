package storetask

import (
	"strconv"
	"strings"
	"testing"

	v "github.com/Cameron-Xie/Golang-API/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidator_Validate(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		input *Task
		err   error
	}{
		{
			input: &Task{},
			err: &v.InvalidParamsError{
				Errs: map[string]string{
					"name": "name is a required field",
				},
			},
		},
		{
			input: &Task{
				Name:        "task_name",
				Description: "task_description",
			},
			err: nil,
		},
		{
			input: &Task{
				Name:        generateLongString(120),
				Description: generateLongString(220),
			},
			err: &v.InvalidParamsError{
				Errs: map[string]string{
					"description": "description must be a maximum of 200 characters in length",
					"name":        "name must be a maximum of 100 characters in length",
				},
			},
		},
	}

	for _, i := range m {
		a.Equal(i.err, NewValidator().Validate(i.input))
	}
}

func generateLongString(size int) string {
	var builder strings.Builder
	for l := range make([]int, size) {
		builder.WriteString(strconv.Itoa(l))
	}

	return builder.String()
}
