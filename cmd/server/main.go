package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/controllers"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Check if required environment variables are set
	requiredEnv := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "PORT"}
	for _, envVar := range requiredEnv {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Required environment variable '%s' is not set", envVar)
		}
	}

	corsMaxAgeHoursStr := os.Getenv("CORS_MAX_AGE_HOURS")
	if corsMaxAgeHoursStr == "" {
		log.Fatalf("Required environment variable 'CORS_MAX_AGE_HOURS' is not set")
	}

	corsMaxAgeHours, err := strconv.Atoi(corsMaxAgeHoursStr)
	if err != nil {
		log.Fatalf("Invalid value for CORS_MAX_AGE_HOURS: %v", err)
	}

	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	router := gin.Default()

	// --- Configuración CORS ---
	// Define las opciones de CORS. Es crucial que el puerto coincida con tu frontend.
	config := cors.DefaultConfig()

	config.AllowOrigins = []string{"http://localhost:5173"} // Origen de tu frontend SvelteKit.
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"} // Necesario si usas Authorization
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true                             // Enviar y recibir cookies session_id
	config.MaxAge = time.Duration(corsMaxAgeHours) * time.Hour // Duración para cachear las respuestas preflight

	// Aplica el middleware CORS a tu router
	router.Use(cors.New(config))
	// --- Fin Configuración CORS ---

	router.POST("/register", controllers.RegisterUser)
	router.POST("/login", controllers.LoginUser)
	router.POST("/logout", controllers.LogoutUser)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", controllers.GetAuthenticatedUser)
		protected.GET("/library", controllers.GetLibrary)
		protected.GET("/audio/:filename", controllers.ServeAudio)
	}
	port := os.Getenv("PORT")
	router.Run(":" + port)
}
