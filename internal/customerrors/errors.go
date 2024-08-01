package customerrors

import "fmt"

// ErrNotFound represents an error for a missing entity with additional context.
type ErrNotFound struct {
	Entity string
	ID     int64
}

// Error implements the error interface for ErrNotFound.
func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s %d not found", e.Entity, e.ID)
}

// NewErrNotFound creates a new ErrNotFound for a given entity and ID.
func NewErrNotFound(entity string, id int64) error {
	return &ErrNotFound{
		Entity: entity,
		ID:     id,
	}
}
