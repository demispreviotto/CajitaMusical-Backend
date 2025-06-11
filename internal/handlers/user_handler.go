package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/dto/user"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/middleware"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// UserHandler es un handler para las operaciones de usuario.
type UserHandler struct {
	userService services.UserServicer // Dependencia del servicio de usuario
}

// NewUserHandler crea una nueva instancia de UserHandler.
func NewUserHandler(userService services.UserServicer) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterUser maneja el registro de nuevos usuarios.
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var userInput user.RegisterUserInput

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResponse, err := h.userService.RegisterUser(context.Background(), userInput)
	if err != nil {
		log.Printf("Handler: User registration failed for %s: %v", userInput.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userResponse)
}

// GetAuthenticatedUser recupera la información del usuario autenticado.
func (h *UserHandler) GetAuthenticatedUser(c *gin.Context) {
	userFromContext, exists := c.Get(middleware.UserContextKey)

	if !exists {
		log.Println("Handler: User info not found in context for /me endpoint (AuthMiddleware might not have run or set it)")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User info not available"})
		return
	}

	userModel, ok := userFromContext.(*models.User)
	if !ok {
		log.Printf("Handler: User info in context is of incorrect type for /me endpoint: got %T", userFromContext)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user info (type mismatch)"})
		return
	}

	safeUser := user.UserInfo{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Name:     userModel.Name,
	}

	c.JSON(http.StatusOK, safeUser)
}

// GetLibrary and ServeAudio would go here too if they are user-specific
// func (h *UserHandler) GetLibrary(c *gin.Context) { ... }
// func (h *UserHandler) ServeAudio(c *gin.Context) { ... }
