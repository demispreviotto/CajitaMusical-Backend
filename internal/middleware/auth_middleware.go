package middleware

import (
	"log"
	"net/http"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db" // Import the db package for interfaces
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserContextKey is the key used to store user information in Gin's context.
const UserContextKey = "user"

// AuthMiddleware struct holds the database dependencies for the middleware.
type AuthMiddleware struct {
	sessionDB db.SessionDBer // Dependency on the SessionDBer interface
	userDB    db.UserDBer    // Dependency on the UserDBer interface
}

// NewAuthMiddleware creates a new instance of AuthMiddleware.
// It takes concrete implementations of SessionDBer and UserDBer.
func NewAuthMiddleware(sessionDB db.SessionDBer, userDB db.UserDBer) *AuthMiddleware {
	return &AuthMiddleware{
		sessionDB: sessionDB,
		userDB:    userDB,
	}
}

// Handler is the actual Gin middleware function.
// It's a method of AuthMiddleware so it can access its dependencies.
func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionIDStr, err := c.Cookie("session_id")
		if err != nil || sessionIDStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session cookie"})
			return
		}

		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			log.Printf("Invalid session ID format: %s. Error: %v", sessionIDStr, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid session ID format"})
			return
		}

		reqCtx := c.Request.Context()

		// Verify the session in the database using the injected sessionDB
		session, err := m.sessionDB.GetSessionByID(reqCtx, sessionID) // <--- FIXED: Using m.sessionDB
		if err != nil {
			log.Printf("Session validation failed for ID %s: %v", sessionID.String(), err)
			// You might want to clear the cookie here if the session is invalid/expired
			c.SetCookie("session_id", "", -1, "/", "localhost", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or expired session"})
			return
		}

		// Get user data associated with the session using the injected userDB
		user, err := m.userDB.GetUserByID(reqCtx, session.UserID)
		if err != nil {
			log.Printf("Failed to retrieve user %s for session %s: %v", session.UserID.String(), sessionID.String(), err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: User data not found"})
			return
		}

		// Store user data in Gin's context
		c.Set(UserContextKey, user)

		// Continue with the next handler in the chain
		c.Next()
	}
}
