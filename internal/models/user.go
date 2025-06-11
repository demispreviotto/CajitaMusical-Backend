package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID        uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username  string    `json:"username" db:"username" gorm:"type:varchar(255);unique;not null"`
	Email     string    `json:"email" db:"email" gorm:"type:varchar(255);unique;not null"`
	Name      string    `json:"name" db:"name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`

	Authentication *Authentication `gorm:"foreignKey:UserID;references:ID"`
}
