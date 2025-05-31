package controllers

import (
	"context"
	"net/http"

	"github.com/demispreviotto/cajitamusical/backend/internal/db"
	"github.com/demispreviotto/cajitamusical/backend/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserInput defines the request body to user registration.
type RegisterUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserResponse defines the response for a successfully registered user.
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

func RegisterUser(c *gin.Context) {
	var userInput RegisterUserInput

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: userInput.Username,
		Email:    userInput.Email,
		Name:     userInput.Name,
		// Password field should likely be handled by db.CreateUser
	}

	if err := db.CreateUser(context.Background(), &user, string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	response := UserResponse{
		ID:       uint(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
	}

	c.JSON(http.StatusCreated, response)
}
