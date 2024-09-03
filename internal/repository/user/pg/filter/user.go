package filter

import "fmt"

// UserFilter is used to filter users by ID or Username (unique fields).
type UserFilter struct {
	ID       *int64
	Username *string
}

// Validate checks that only one field (ID or Username) is set in the UserFilter.
// It returns an error if both or neither are set.
func (f *UserFilter) Validate() error {
	if f.ID != nil && f.Username != nil {
		return fmt.Errorf("only one of ID or Username should be set")
	}
	if f.ID == nil && f.Username == nil {
		return fmt.Errorf("either ID or Username must be set")
	}
	return nil
}
