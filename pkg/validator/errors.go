package validator

import (
	"fmt"
	"strings"
)

type InvalidParamsError struct {
	Errs map[string]string
}

func (e *InvalidParamsError) Error() string {
	l := make([]string, 0)
	for k, v := range e.Errs {
		l = append(l, fmt.Sprintf("%s: %s", k, v))
	}

	return strings.Join(l, ",")
}
