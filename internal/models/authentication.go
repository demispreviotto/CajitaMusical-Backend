package models

import (
	"time"

	"github.com/google/uuid"
)

// Authentication stores user authentication details.
type Authentication struct {
	UserID       uuid.UUID  `json:"user_id" db:"user_id" gorm:"primaryKey;type:uuid;not null"`
	PasswordHash string     `json:"password_hash" db:"password_hash" gorm:"type:varchar(255);not null"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at" gorm:"not null"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at" gorm:"not null"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}
