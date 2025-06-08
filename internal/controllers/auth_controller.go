package controllers

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
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
		log.Printf("Login failed for username %s: %v", loginRequest.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginRequest.Password)); err != nil {
		log.Printf("Login failed for username %s: password mismatch", loginRequest.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create a new session
	sessionID := uuid.New().String()

	sessionDurationStr := os.Getenv("SESSION_DURATION_HOURS")
	sessionDurationHours := 24 // Default
	if dur, err := time.ParseDuration(sessionDurationStr + "h"); err == nil {
		sessionDurationHours = int(dur.Hours())
	}
	expiresAt := time.Now().Add(time.Duration(sessionDurationHours) * time.Hour)

	// Obtener la IP del cliente como string
	clientIPStr := c.ClientIP()
	// Convertir la string de IP a net.IP
	clientIP := net.ParseIP(clientIPStr)

	if clientIP == nil {
		log.Printf("Warning: Could not parse client IP '%s' from Gin. Using nil/default.", clientIPStr)
	}

	session := &models.Session{
		SessionID: sessionID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		UserAgent: c.Request.UserAgent(),
		IPAddress: clientIP,
	}

	if err := db.CreateSession(context.Background(), session); err != nil {
		log.Printf("Failed to create session for user %d: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Set the session cookie
	c.SetCookie("session_id", sessionID, int(time.Until(expiresAt).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, LoginResponse{Message: "Login successful", User: UserInfo{ID: uint(user.ID), Username: user.Username, Email: user.Email, Name: user.Name}})
}

// LogoutUser handles user logout by invalidating the session.
func LogoutUser(c *gin.Context) {
	// Intenta obtener el ID de la sesión de la cookie
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Already logged out or session expired."})
		return
	}

	if err := db.DeleteSession(context.Background(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout fully, please try again."})
		return
	}

	// Invalida la cookie del navegador para que el frontend ya no la envíe.
	c.SetCookie("session_id", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func CleanupExpiredSessions(c *gin.Context) {
	log.Println("Manual cleanup of expired sessions requested...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.DeleteExpiredSessions(ctx)
	if err != nil {
		log.Printf("Error during manual session cleanup: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup expired sessions"})
		return
	}

	log.Println("Manual expired session cleanup completed.")
	c.JSON(http.StatusOK, gin.H{"message": "Expired sessions cleanup initiated successfully"})
}
