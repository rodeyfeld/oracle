package order

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func TestConnectPostgres(t *testing.T) {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := &PostgresDB{}
	db.Connect()
}
