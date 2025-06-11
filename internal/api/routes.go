package api

import (
	"os"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/handlers"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/middleware"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes for the application.
func SetupRoutes(router *gin.Engine) {
	// Initialize DB layer implementations
	userDB := db.NewUserDB()
	sessionDB := db.NewSessionDB()
	songDB := db.NewSongDB()

	// Get environment variable for music directory
	musicDir := os.Getenv("MUSIC_DIRECTORY")
	if musicDir == "" {
		panic("MUSIC_DIRECTORY environment variable is not set!")
	}

	// Initialize service layer implementations
	userService := services.NewUserService(userDB)
	authService := services.NewAuthService(userDB, sessionDB)
	songService := services.NewSongService(songDB, musicDir)

	// Initialize handlers with their service dependencies
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	songHandler := handlers.NewSongHandler(songService)

	// Initialize the AuthMiddleware with its DB dependencies
	authMiddleware := middleware.NewAuthMiddleware(sessionDB, userDB)

	// Public routes (no authentication required)
	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", authHandler.LoginUser)

	// Protected routes (require AuthMiddleware)
	protected := router.Group("/")
	protected.Use(authMiddleware.Handler())
	{
		protected.GET("/me", userHandler.GetAuthenticatedUser)
		protected.POST("/logout", authHandler.LogoutUser)

		// Song routes
		protected.GET("/library", songHandler.GetLibrary)
		protected.GET("/audio/:filename", songHandler.ServeAudio)
	}

	// Admin routes
	admin := router.Group("/admin")
	admin.Use(authMiddleware.Handler())
	{
		admin.POST("/cleanup-sessions", authHandler.CleanupExpiredSessions)
	}
}
