package auth

// LoginUserInput defines the request body for user login.
type LoginUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
