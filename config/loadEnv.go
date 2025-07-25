package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
		log.Println("External set up of environmental variables at production mode")
	}
}
