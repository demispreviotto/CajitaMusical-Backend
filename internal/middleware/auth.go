package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// You might need to import your db package here if you plan to verify sessions
)

// AuthMiddleware is a basic authentication middleware.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not logged in"})
			return
		}
		// In a real application, you would verify the session ID against the database here
		c.Next()
	}
}
