package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Declare config variables
var MONGO_URI string
var JWT_SECRET string

// init runs automatically when the package is imported
func init() {
	// Load .env file from project root
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	// Read environment variables
	MONGO_URI = os.Getenv("MONGO_URI")
	JWT_SECRET = os.Getenv("JWT_SECRET")

	// Optional: Warn if values are missing
	if MONGO_URI == "" {
		log.Println("Warning: MONGO_URI is not set!")
	}
	if JWT_SECRET == "" {
		log.Println("Warning: JWT_SECRET is not set!")
	}
}
