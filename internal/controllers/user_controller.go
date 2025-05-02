package controllers

import (
	"context"
	"net/http"

	"github.com/demispreviotto/cajitamusical/backend/internal/db"     // Replace with your module path
	"github.com/demispreviotto/cajitamusical/backend/internal/models" // Replace with your module path
	"github.com/gin-gonic/gin"
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
