package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB es la instancia global de la conexión a la base de datos GORM.
var DB *gorm.DB

// Connect inicializa la conexión a la base de datos PostgreSQL usando GORM.
func Connect() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/New_York",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Loggea todas las consultas de info
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully with GORM!")
	return nil
}

// Close cierra la conexión a la base de datos.
func Close() {
	if DB == nil { // Asegurarse de que DB no sea nil antes de intentar cerrar
		return
	}
	sqlDB, err := DB.DB() // GORM te da el *sql.DB subyacente para cerrar
	if err != nil {
		log.Printf("Error getting underlying DB connection: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
	log.Println("Database connection closed.")
}

// GetDB returns the GORM database instance.
func GetDB() *gorm.DB {
	return DB
}
