package db

import (
	"context"
	"fmt"

	"github.com/demispreviotto/cajitamusical/backend/internal/models" // Replace with your module path
)

// CreateUser creates a new user in the database.
func CreateUser(ctx context.Context, user *models.User, passwordHash string) error {
	_, err := DB.Exec(ctx,
		"INSERT INTO users (username, email, name) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at",
		user.Username, user.Email, user.Name,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	_, err = DB.Exec(ctx,
		"INSERT INTO authentication (user_id, password_hash) VALUES ((SELECT id FROM users WHERE username = $1), $2)",
		user.Username, passwordHash,
	)
	if err != nil {
		return fmt.Errorf("failed to insert authentication details: %w", err)
	}

	return nil
}

// GetUserByUsername retrieves a user from the database by their username.
func GetUserByUsername(ctx context.Context, username string) (*models.User, string, error) {
	user := &models.User{}
	var passwordHash string
	err := DB.QueryRow(ctx,
		"SELECT u.id, u.username, u.email, u.name, u.created_at, u.updated_at, a.password_hash "+
			"FROM users u JOIN authentication a ON u.id = a.user_id "+
			"WHERE u.username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt, &passwordHash)

	if err != nil {
		return nil, "", fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, passwordHash, nil
}
