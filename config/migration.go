package config

import (
	"flag"
	"food-delivery-app-server/infrastructure"
	"log"
	"os"
)

func HandleMigrationFlag() {
	migrateFlag := flag.Bool("migrate", false, "run db migration and exit")
	flag.Parse()

	if *migrateFlag {
		infrastructure.ConnectDb()

		log.Println("Starting automigrate...")
		infrastructure.SyncDatabase()

		log.Println("Migration complete... exiting...")
		os.Exit(0)
	} else {
		log.Println("Skipping db migation (flag was not set)")
	}
}
