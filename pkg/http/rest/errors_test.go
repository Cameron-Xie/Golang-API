package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpError_Error(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		err      *HTTPError
		expected string
	}{
		{
			err: &HTTPError{
				Title: "title",
				InvalidParams: []InvalidParam{
					{
						Name:   "id",
						Reason: "invalid id",
					},
				},
				StatusCode: 400,
			},
			expected: `{"title":"title","invalid_params":[{"name":"id","reason":"invalid id"}]}`,
		},
		{
			err: &HTTPError{
				Title:      "title",
				StatusCode: 400,
			},
			expected: `{"title":"title"}`,
		},
	}

	for _, i := range m {
		a.Equal(i.expected, i.err.Error())
	}
}
