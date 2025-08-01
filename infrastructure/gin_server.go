package infrastructure

import (
	"food-delivery-app-server/infrastructure/routes"
	"food-delivery-app-server/middleware"

	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunGin(corsConfig cors.Config) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.Use(cors.New(corsConfig))

	if os.Getenv("GIN_MODE") == "release" {
		err := r.SetTrustedProxies([]string{"127.0.0.1"})
		if err != nil {
			log.Fatalf("Failed to set trusted proxies: %v", err)
		}
	}

	r.Use(middleware.ErrorHandler())

	routes.RegisterRoutes(r, DB, RedisClient)

	err := r.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatal("Failed to run the server")
	}

}
