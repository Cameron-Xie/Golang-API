package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError_Error(t *testing.T) {
	a := assert.New(t)
	m := []struct {
		err      NotFoundError
		excepted string
	}{
		{
			err: NotFoundError{
				Table: "table_name",
				Value: "id",
			},
			excepted: `id is not found in table_name`,
		},
	}

	for _, i := range m {
		a.Equal(i.excepted, i.err.Error())
	}
}
