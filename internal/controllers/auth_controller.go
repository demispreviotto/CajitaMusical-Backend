package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/demispreviotto/cajitamusical/backend/internal/db"
	"github.com/demispreviotto/cajitamusical/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// LoginUserInput defines the request body for user login.
type LoginUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserInfo defines the user information in the login response.
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

// LoginResponse defines the response for a successful login.
type LoginResponse struct {
	Message string   `json:"message"`
	User    UserInfo `json:"user"`
}

func LoginUser(c *gin.Context) {
	var loginRequest LoginUserInput

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

	c.JSON(http.StatusOK, LoginResponse{Message: "Login successful", User: UserInfo{ID: uint(user.ID), Username: user.Username, Email: user.Email, Name: user.Name}})
}
