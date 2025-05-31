package main

import (
	"log"
	"os"

	"github.com/demispreviotto/cajitamusical/backend/internal/controllers"
	"github.com/demispreviotto/cajitamusical/backend/internal/db"
	"github.com/demispreviotto/cajitamusical/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Cajita Musical API
// @version 1.0
// @description This is the API for the Cajita Musical application.

// @contact.name Demis Previotto
// @contact.url http://demispreviotto.com
// @contact.email demis.previotto@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

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

	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	router := gin.Default()

	router.POST("/register", controllers.RegisterUser)
	router.POST("/login", controllers.LoginUser)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/library", controllers.GetLibrary)
		protected.GET("/audio/:filename", controllers.ServeAudio)
	}
	port := os.Getenv("PORT")
	router.Run(":" + port)
}
