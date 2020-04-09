package storetask

import (
	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required,max=100"`
	Description string    `json:"description" validate:"max=200"`
}
