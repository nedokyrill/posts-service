package app

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/nedokyrill/posts-service/pkg/db"
	"github.com/nedokyrill/posts-service/pkg/logger"
)

func Run() {
	logger.InitLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatal("Error loading .env file")
	}

	if os.Getenv("IN_MEM_STORAGE") == "true" {
	} else {

	}

	conn, err := db.Connect()
	if err != nil {
		logger.Logger.Fatal(err)
	}
	defer conn.Close()
}
