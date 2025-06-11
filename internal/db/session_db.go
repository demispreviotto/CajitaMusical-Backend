package db

import (
	"context"
	"time"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
	"github.com/google/uuid"
)

// SessionDBer define la interfaz para las operaciones de la base de datos de sesiones.
type SessionDBer interface {
	CreateSession(ctx context.Context, session *models.Session) error
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	DeleteExpiredSessions(ctx context.Context) error
	// Agrega otros métodos de DB de sesión aquí
}

// sessionDB es la implementación concreta de SessionDBer.
type sessionDB struct{}

// NewSessionDB crea una nueva instancia de SessionDB.
func NewSessionDB() SessionDBer {
	return &sessionDB{}
}

// Implementación de CreateSession
func (sdb *sessionDB) CreateSession(ctx context.Context, session *models.Session) error {
	return DB.WithContext(ctx).Create(session).Error
}

// Implementación de GetSessionByID
func (sdb *sessionDB) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := DB.WithContext(ctx).Where("session_id = ? AND expires_at > ?", sessionID, time.Now()).First(&session).Error
	return &session, err
}

// Implementación de DeleteSession
func (sdb *sessionDB) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	return DB.WithContext(ctx).Delete(&models.Session{}, "session_id = ?", sessionID).Error
}

// Implementación de DeleteExpiredSessions
func (sdb *sessionDB) DeleteExpiredSessions(ctx context.Context) error {
	return DB.WithContext(ctx).Where("expires_at <= ?", time.Now()).Delete(&models.Session{}).Error
}
