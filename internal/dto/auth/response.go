package auth

import (
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/dto/user"
)

// LoginResponse defines the response for a successful login.
type LoginResponse struct {
	Message string        `json:"message"`
	User    user.UserInfo `json:"user"`
}
