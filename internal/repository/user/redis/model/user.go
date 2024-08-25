package model

// User represents a user entity in the Redis database.
type User struct {
	ID          int64  `redis:"id"`
	Username    string `redis:"username"`
	Email       string `redis:"email"`
	Role        string `redis:"role"`
	CreatedAtNs int64  `redis:"created_at"`
	UpdatedAtNs int64  `redis:"updated_at"`
}
