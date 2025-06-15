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
	songService := services.NewSongService(songDB)

	// Initialize handlers with their service dependencies
	userHandler := handlers.NewuserHandler(userService)
	authHandler := handlers.NewauthHandler(authService)
	songHandler := handlers.NewsongHandler(songService)

	// Initialize the AuthMiddleware with its DB dependencies
	authMiddleware := middleware.NewAuthMiddleware(sessionDB, userDB)

	// Public routes (no authentication required)
	public := router.Group("/api")
	{
		public.POST("/register", userHandler.RegisterUser)
		public.POST("/login", authHandler.LoginUser)
	}
	// Protected routes (require AuthMiddleware)
	protected := router.Group("/api")
	protected.Use(authMiddleware.Handler())
	{
		protected.GET("/me", userHandler.GetAuthenticatedUser)
		protected.POST("/logout", authHandler.LogoutUser)

		// Song routes
		protected.GET("/library", songHandler.GetLibrary)
		protected.GET("/audio/:filename", songHandler.ServeAudio)
		protected.GET("/album-art/*filepath", songHandler.ServeAlbumArt)
	}
	// Admin routes
	admin := router.Group("/api/admin")
	admin.Use(authMiddleware.Handler())
	{
		admin.POST("/cleanup-sessions", authHandler.CleanupExpiredSessions)
		admin.POST("/scan-music", songHandler.ScanMusicLibrary)
	}
}
