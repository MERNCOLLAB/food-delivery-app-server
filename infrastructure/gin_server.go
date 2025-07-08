package infrastructure

import (
	"food-delivery-app-server/infrastructure/routes"
	"food-delivery-app-server/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func RunGin() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.Use(middleware.ErrorHandler())

	routes.RegisterRoutes(r, DB, RedisClient)

	err := r.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to run the server")
	}

}
