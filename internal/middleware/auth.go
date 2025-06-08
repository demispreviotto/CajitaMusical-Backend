package middleware

import (
	"log"
	"net/http"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db"
	"github.com/gin-gonic/gin"
)

// UserContextKey es la clave utilizada para almacenar la información del usuario en el contexto de Gin.
const UserContextKey = "user"

// AuthMiddleware es un middleware de autenticación que valida la sesión y adjunta la información del usuario.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session cookie"})
			return
		}

		reqCtx := c.Request.Context() // Usamos el contexto de la solicitud

		// 1. Verificar la sesión en la base de datos
		session, err := db.GetSessionByID(reqCtx, sessionID)
		if err != nil {
			log.Printf("Session validation failed for ID %s: %v", sessionID, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or expired session"})
			return
		}

		// 2. Obtener los datos del usuario asociados a la sesión
		user, err := db.GetUserByID(reqCtx, session.UserID)
		if err != nil {
			log.Printf("Failed to retrieve user %d for session %s: %v", session.UserID, sessionID, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: User data not found"})
			return
		}

		// 3. Almacenar los datos del usuario en el contexto de Gin
		c.Set(UserContextKey, user)

		// Continuar con la siguiente cadena de manejadores
		c.Next()
	}
}
