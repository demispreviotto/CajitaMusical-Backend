package models

import (
	"net"
	"time"

	"github.com/google/uuid"
)

// Session represents an active user session.
type Session struct {
	SessionID uuid.UUID `json:"session_id" db:"session_id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" db:"user_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at" gorm:"not null"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	IPAddress net.IP    `json:"ip_address" db:"ip_address"`
}
