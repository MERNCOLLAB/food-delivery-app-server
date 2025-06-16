package infrastructure

import (
	"food-delivery-app-server/models"
	"log"
)

func SyncDatabase() {
	log.Println("Syncing declared database schema...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Restaurant{},
		&models.Address{},
		&models.MenuItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Notification{},
		&models.PasswordReset{},
	)

	if err != nil {
		log.Fatalf("Failed to auto-migrate the provided db schema: %v", err)
		return
	}

	log.Println("Database migration is successful")
}
