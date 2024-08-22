package validators

import (
	"errors"
)

// ValidatePassword provides simple password validation.
func ValidatePassword(password, passwordConfirm string) error {
	if password != passwordConfirm {
		return errors.New("passwords don't match")
	}
	return nil
}
