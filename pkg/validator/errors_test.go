package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidParam_Error(t *testing.T) {
	a := assert.New(t)
	e := &InvalidParamsError{
		Errs: map[string]string{
			"error":  "something went wrong",
			"error2": "another error",
		},
	}

	a.Contains(e.Error(), "error: something went wrong")
	a.Contains(e.Error(), "error2: another error")
}
