package model

import (
	"time"
)

// User represents a business logic user model.
type User struct {
	ID        int64
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
