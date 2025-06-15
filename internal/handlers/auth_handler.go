package handlers

import (
	"context"
	"log"
	"net/http" // Necesario para net.ParseIP
	"time"     // Necesario para time.Now()

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/dto/auth"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/services" // Importar el paquete services
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// authHandler es un handler para las operaciones de autenticación.
type authHandler struct {
	authService services.AuthServicer // Dependencia del servicio de autenticación
}

// NewauthHandler crea una nueva instancia de authHandler.
func NewauthHandler(authService services.AuthServicer) *authHandler {
	return &authHandler{authService: authService}
}

// LoginUser maneja el inicio de sesión del usuario.
func (h *authHandler) LoginUser(c *gin.Context) {
	var loginRequest auth.LoginUserInput

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := c.Request.UserAgent()
	clientIP := c.ClientIP()

	loginResponse, session, err := h.authService.Login(context.Background(), loginRequest.Username, loginRequest.Password, userAgent, clientIP)
	if err != nil {
		log.Printf("Handler: Login attempt failed for %s: %v", loginRequest.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}) // Retorna el error del servicio
		return
	}

	// Set the session cookie
	// Calcular el tiempo de vida de la cookie basado en la expiración de la sesión
	maxAge := int(time.Until(session.ExpiresAt).Seconds())
	c.SetCookie("session_id", session.SessionID.String(), maxAge, "/", "localhost", false, true)

	c.JSON(http.StatusOK, loginResponse)
}

// LogoutUser maneja el cierre de sesión del usuario.
func (h *authHandler) LogoutUser(c *gin.Context) {
	sessionIDStr, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Already logged out or session expired."})
		return
	}

	sessionUUID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		log.Printf("Handler: Invalid session ID format from cookie: %s. Error: %v", sessionIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	if err := h.authService.Logout(context.Background(), sessionUUID); err != nil {
		log.Printf("Handler: Failed to logout session %s: %v", sessionUUID.String(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Retorna el error del servicio
		return
	}

	c.SetCookie("session_id", "", -1, "/", "localhost", false, true) // Invalida la cookie
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// CleanupExpiredSessions maneja la limpieza manual de sesiones expiradas (podría ser un endpoint de admin).
func (h *authHandler) CleanupExpiredSessions(c *gin.Context) {
	log.Println("Handler: Manual cleanup of expired sessions requested...")

	if err := h.authService.CleanupExpiredSessions(context.Background()); err != nil {
		log.Printf("Handler: Error during manual session cleanup: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Retorna el error del servicio
		return
	}

	log.Println("Handler: Manual expired session cleanup completed.")
	c.JSON(http.StatusOK, gin.H{"message": "Expired sessions cleanup initiated successfully"})
}
