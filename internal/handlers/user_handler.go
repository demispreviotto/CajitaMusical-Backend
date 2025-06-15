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

// userHandler es un handler para las operaciones de usuario.
type userHandler struct {
	userService services.UserServicer // Dependencia del servicio de usuario
}

// NewuserHandler crea una nueva instancia de userHandler.
func NewuserHandler(userService services.UserServicer) *userHandler {
	return &userHandler{userService: userService}
}

// RegisterUser maneja el registro de nuevos usuarios.
func (h *userHandler) RegisterUser(c *gin.Context) {
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
func (h *userHandler) GetAuthenticatedUser(c *gin.Context) {
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
// func (h *userHandler) GetLibrary(c *gin.Context) { ... }
// func (h *userHandler) ServeAudio(c *gin.Context) { ... }
