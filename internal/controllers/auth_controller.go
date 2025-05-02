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

// LoginUser handles user login.
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
