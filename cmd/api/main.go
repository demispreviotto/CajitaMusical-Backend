package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/api"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
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
	requiredEnv := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "PORT", "MUSIC_DIRECTORY", "FRONTEND_ORIGIN"}
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

	// AutoMigrar tus modelos GORM
	log.Println("Running GORM auto-migrations...")
	err = db.DB.AutoMigrate(
		&models.User{},
		&models.Authentication{},
		&models.Session{},
		&models.Song{},
		// Añade cualquier otro modelo que GORM deba gestionar aquí
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("GORM auto-migrations completed.")

	routerEngine := gin.Default()

	// --- Configuración CORS ---
	config := cors.DefaultConfig()

	// Get frontend origin from environment variable
	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		// This check is technically redundant if "FRONTEND_ORIGIN" is in requiredEnv,
		// but it adds an extra layer of safety or clearer error message if requiredEnv is changed.
		log.Fatalf("FRONTEND_ORIGIN environment variable is not set.")
	}

	config.AllowOrigins = []string{frontendOrigin} // Origen de tu frontend SvelteKit.
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"} // Necesario si usas Authorization
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true                             // Enviar y recibir cookies session_id
	config.MaxAge = time.Duration(corsMaxAgeHours) * time.Hour // Duración para cachear las respuestas preflight

	// Aplica el middleware CORS a tu router
	routerEngine.Use(cors.New(config))
	// --- Fin Configuración CORS ---

	api.SetupRoutes(routerEngine)

	port := os.Getenv("PORT")
	routerEngine.Run(":" + port)
}
