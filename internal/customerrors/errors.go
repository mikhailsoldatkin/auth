package customerrors

import "fmt"

// ErrNotFound represents an error for a missing entity with additional context.
type ErrNotFound struct {
	Entity     string
	Identifier string
}

// Error implements the error interface for ErrNotFound.
func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s with %s not found", e.Entity, e.Identifier)
}

// NewErrNotFound creates a new ErrNotFound for a given entity and identifier.
func NewErrNotFound(entity string, identifier any) error {
	var idStr string
	switch v := identifier.(type) {
	case int64:
		idStr = fmt.Sprintf("ID %d", v)
	case string:
		idStr = fmt.Sprintf("username '%s'", v)
	default:
		idStr = fmt.Sprintf("unknown identifier %v", v)
	}
	return &ErrNotFound{
		Entity:     entity,
		Identifier: idStr,
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
