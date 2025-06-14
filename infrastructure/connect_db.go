package infrastructure

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() {
	dbUrl := os.Getenv("DATABASE_URL")
	dbName := os.Getenv("DB_NAME")

	if dbUrl == "" {
		log.Fatal("Environmental variable DATABASE_URL does not exist")
	}

	if dbName == "" {
		log.Fatal("Unkown DB_NAME")
	}

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("%v is connected successfully", err)
	}

	DB = db

	log.Printf("🔗 Connected successfully to %s database", dbName)

}
