package main

import (
	"food-delivery-app-server/config"
	"food-delivery-app-server/infrastructure"
)

func main() {
	config.LoadEnvVariables()

	infrastructure.ConnectDb()

	infrastructure.RunGin()
}
