package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"

	"github.com/sgomeza13/stock-recommender/utils"
)

// DB is the global database connection
var DB *pgx.Conn

// ConnectDB initializes the database connection
func ConnectDB() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Read DB credentials from .env

	dsn := utils.GetDSN(false)

	// Connect to the database
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Assign the connection to the global variable
	DB = conn
	log.Println("Successfully connected to the database!")
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close(context.Background())
		log.Println("Database connection closed.")
	}
}
