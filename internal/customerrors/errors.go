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

// ErrInvalidPassword represents an error when the provided password is incorrect.
type ErrInvalidPassword struct{}

// Error implements the error interface for ErrInvalidPassword.
func (e *ErrInvalidPassword) Error() string {
	return "invalid password"
}

// NewErrInvalidPassword creates a new ErrInvalidPassword.
func NewErrInvalidPassword() error {
	return &ErrInvalidPassword{}
}

// ErrInvalidToken represents an error when the provided token is invalid.
type ErrInvalidToken struct{}

// Error implements the error interface for ErrInvalidPassword.
func (e *ErrInvalidToken) Error() string {
	return "invalid token"
}

// NewErrInvalidToken creates a new NewErrInvalidToken.
func NewErrInvalidToken() error {
	return &ErrInvalidToken{}
}

// ErrForbidden represents an error for a forbidden action.
type ErrForbidden struct{}

// Error implements the error interface for ErrForbidden.
func (e *ErrForbidden) Error() string {
	return "access denied"
}

// NewErrForbidden creates a new ErrForbidden error.
func NewErrForbidden() error {
	return &ErrForbidden{}
}
