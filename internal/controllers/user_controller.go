package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/middleware"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
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
func GetAuthenticatedUser(c *gin.Context) {
	// El middleware AuthMiddleware debería haber puesto la información del usuario en el contexto
	// usando la clave definida en middleware.UserContextKey
	user, exists := c.Get(middleware.UserContextKey)

	if !exists {
		log.Println("User info not found in context for /me endpoint (AuthMiddleware might not have run or set it)")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User info not available"})
		return
	}

	// Asegúrate de que el tipo de 'user' coincida con lo que tu middleware pone en el contexto.
	// Si tu middleware pone directamente el *models.User, entonces sería:
	userModel, ok := user.(*models.User) // <- Nota el puntero si tu middleware lo devuelve así
	if !ok {
		// Fallback si no es *models.User, tal vez sea models.User directamente o UserInfo
		// Esto depende de cómo tu AuthMiddleware esté estructurado.
		// Si tienes un UserInfo struct para la respuesta, puedes intentar un type assertion a eso también.
		log.Printf("User info in context is of incorrect type for /me endpoint: got %T", user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user info (type mismatch)"})
		return
	}

	// Si tu modelo User incluye campos sensibles, crea una versión 'segura' o un DTO para la respuesta.
	// Usaremos UserInfo que ya tienes definida para el LoginResponse
	safeUser := UserInfo{
		ID:       uint(userModel.ID),
		Username: userModel.Username,
		Email:    userModel.Email,
		Name:     userModel.Name,
	}

	c.JSON(http.StatusOK, safeUser)
}
