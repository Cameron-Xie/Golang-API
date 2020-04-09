package postgres

import "fmt"

type NotFoundError struct {
	Table string
	Value string
}

func (err *NotFoundError) Error() string {
	return fmt.Sprintf("%v is not found in %v", err.Value, err.Table)
}
