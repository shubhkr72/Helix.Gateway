package auth

import (
	"context"
	"database/sql"
)

// PostgresRepository implements the Repository interface using PostgreSQL.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository.
func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{
		db: db,
	}
}

// CreateUser inserts a new user into the users table.
func (r *PostgresRepository) CreateUser(ctx context.Context, user *User) error {

	query := `
	INSERT INTO users
	(
		id,
		email,
		password_hash,
		role,
		created_at,
		updated_at
	)
	VALUES
	(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6
	)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// GetByEmail fetches a user by email.
func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {

	query := `
	SELECT
		id,
		email,
		password_hash,
		role,
		created_at,
		updated_at
	FROM users
	WHERE email = $1
	`

	user := &User{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}