package db

import (
	"context"
	"fmt"

	"github.com/demispreviotto/cajitamusical/backend/internal/models"
)

// CreateSession creates a new session in the database.
func CreateSession(ctx context.Context, session *models.Session) error {
	_, err := DB.Exec(ctx,
		"INSERT INTO sessions (session_id, user_id, created_at, expires_at, user_agent, ip_address) "+
			"VALUES ($1, $2, $3, $4, $5, $6)",
		session.SessionID, session.UserID, session.CreatedAt, session.ExpiresAt, session.UserAgent, session.IPAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}
	return nil
}
