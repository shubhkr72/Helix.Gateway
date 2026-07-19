package auth

import "context"

// Repository defines all database operations required by the auth service.
type Repository interface {
	// CreateUser inserts a new user into the database.
	CreateUser(ctx context.Context, user *User) error

	// GetByEmail retrieves a user by email.
	GetByEmail(ctx context.Context, email string) (*User, error)
}