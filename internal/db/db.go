package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// DB holds the PostgreSQL database connection pool.
var DB *pgx.Conn

// Connect establishes a connection to the PostgreSQL database.
func Connect() error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = conn
	fmt.Println("Successfully connected to PostgreSQL!")
	return nil
}

// Close closes the database connection pool.
func Close() error {
	if DB != nil {
		return DB.Close(context.Background())
	}
	return nil
}
