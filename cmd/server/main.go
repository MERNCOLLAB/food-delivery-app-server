package main

import (
	"food-delivery-app-server/config"
	"food-delivery-app-server/infrastructure"
)

func main() {
	config.LoadEnvVariables()

	// add --migrate in running Go if it needs db migration
	config.HandleMigrationFlag()

	infrastructure.ConnectRedis()
	infrastructure.ConnectDb()
	infrastructure.RunGin()
}
