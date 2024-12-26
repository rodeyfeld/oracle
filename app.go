package application

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"oracle/internal/scholar"
)

var DebugMode bool

func Start() {

	// Env handling
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	debugModeStr := os.Getenv("DEBUG_MODE")

	DebugMode = false
	if debugModeStr == "true" {
		DebugMode = true
	}
	log.Printf("Debug mode=%t", DebugMode)
	// Scrape
	scholar.Enscribe()
}
