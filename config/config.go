package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
