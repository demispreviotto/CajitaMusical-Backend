package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB holds the PostgreSQL database connection pool.
var DB *pgxpool.Pool

// Connect establishes a connection to the PostgreSQL database.
func Connect() error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = pool
	fmt.Println("Successfully connected to PostgreSQL!")
	return nil
}

// Close closes the database connection pool.
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// GetDB returns the database connection pool.
func GetDB() *pgxpool.Pool {
	return DB
}
