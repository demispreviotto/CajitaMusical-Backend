package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/demispreviotto/cajitamusical/backend/internal/db"     // Replace with your module path
	"github.com/demispreviotto/cajitamusical/backend/internal/models" // Replace with your module path
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles the registration of a new user.
func RegisterUser(c *gin.Context) {
	var user models.User
	var password string

	if err := c.ShouldBindJSON(&struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}{
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
		Password: password,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create the user in the database
	if err := db.CreateUser(context.Background(), &user, string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
func LoginUser(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, hashedPassword, err := db.GetUserByUsername(context.Background(), loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create a new session
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour) // Session expires in 24 hours
	session := &models.Session{
		SessionID: sessionID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		UserAgent: c.Request.UserAgent(),
		IPAddress: c.ClientIP(),
	}

	if err := db.CreateSession(context.Background(), session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Set the session cookie
	c.SetCookie("session_id", sessionID, int(time.Until(expiresAt).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": gin.H{"id": user.ID, "username": user.Username, "email": user.Email, "name": user.Name}})
}
