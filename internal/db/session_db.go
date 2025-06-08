package db

import (
	"context"
	"fmt"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
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

// GetSessionByID busca una sesión por su ID y verifica si está expirada.
func GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	var session models.Session
	// Aquí, la consulta busca la sesión y también verifica si no ha expirado aún.
	// La cláusula WHERE `expires_at > NOW()` es crucial.
	err := DB.QueryRow(ctx,
		"SELECT session_id, user_id, created_at, expires_at, user_agent, ip_address FROM sessions WHERE session_id = $1 AND expires_at > NOW()",
		sessionID,
	).Scan(
		&session.SessionID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
		&session.UserAgent,
		&session.IPAddress,
	)
	if err != nil {
		// sql.ErrNoRows significa que la sesión no se encontró o está expirada
		return nil, fmt.Errorf("failed to get session by ID: %w", err)
	}
	return &session, nil
}

// DeleteSession elimina una sesión específica por su SessionID.
// Esto se usará para el logout explícito.
func DeleteSession(ctx context.Context, sessionID string) error {
	result, err := DB.Exec(ctx, "DELETE FROM sessions WHERE session_id = $1", sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session %s: %w", sessionID, err)
	}
	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("session %s not found for deletion", sessionID)
	}
	return nil
}

// DeleteExpiredSessions elimina las sesiones caducadas de la base de datos.
// Esto se usará para la "recolección de basura" periódica.
func DeleteExpiredSessions(ctx context.Context) error {
	result, err := DB.Exec(ctx, "DELETE FROM sessions WHERE expires_at < NOW()")
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}
	rowsAffected := result.RowsAffected()
	// The _ = rowsAffected is fine if you're not using rowsAffected in this function.
	_ = rowsAffected
	return nil
}
