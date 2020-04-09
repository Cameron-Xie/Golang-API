package updatetask

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	val "github.com/Cameron-Xie/Golang-API/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidator_Validate(t *testing.T) {
	a := assert.New(t)
	v := NewValidator()
	m := []struct {
		input    map[string]interface{}
		expected map[string]interface{}
		err      error
	}{
		{
			input: map[string]interface{}{},
			err:   errors.New("empty input"),
		},
		{
			input: map[string]interface{}{
				"name":        "task_name",
				"description": "description",
			},
			expected: map[string]interface{}{
				"name":        "task_name",
				"description": "description",
			},
		},
		{
			input: map[string]interface{}{
				"description": "description",
			},
			expected: map[string]interface{}{
				"description": "description",
			},
		},
		{
			input: map[string]interface{}{
				"name": "task_name",
			},
			expected: map[string]interface{}{
				"name": "task_name",
			},
		},
		{
			input: map[string]interface{}{
				"name":        generateLongString(120),
				"description": generateLongString(220),
			},
			err: &val.InvalidParamsError{
				Errs: map[string]string{
					"description": "description must not be longer than 200 characters",
					"name":        "name must be at least 1 character and less than 100 characters",
				},
			},
		},
	}

	for _, i := range m {
		res, err := v.Validate(i.input)
		a.Equal(i.expected, res)
		a.Equal(i.err, err)
	}
}

func generateLongString(size int) string {
	var builder strings.Builder
	for l := range make([]int, size) {
		builder.WriteString(strconv.Itoa(l))
	}

	return builder.String()
}
