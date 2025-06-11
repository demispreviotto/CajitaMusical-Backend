package user

import "github.com/google/uuid" // Import uuid for the ID field

// UserInfo defines the user information to be exposed in responses.
type UserInfo struct {
	ID       uuid.UUID `json:"id"` // Changed to uuid.UUID
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
}

// // LoginResponse defines the response for a successful login.
// type LoginResponse struct {
// 	Message string   `json:"message"`
// 	User    UserInfo `json:"user"`
// }

// You might also move UserResponse here if it's purely for API responses
type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
}
