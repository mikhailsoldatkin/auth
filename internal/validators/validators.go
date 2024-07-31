package validators

import (
	"errors"
	"regexp"
)

const (
	emailRegex        = `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	passwordMinLength = 8
)

// ValidateEmail checks if the given email address is in valid format.
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ValidatePassword provides simple password validation.
func ValidatePassword(password, passwordConfirm string) error {
	if len(password) < passwordMinLength {
		return errors.New("password must be at least 8 characters long")
	}
	if password != passwordConfirm {
		return errors.New("passwords don't match")
	}
	return nil
}
