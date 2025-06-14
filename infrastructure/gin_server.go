package infrastructure

import (
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

	RegisterRoutes(r)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is up and running!",
		})
	})

	err := r.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to run the server")
	}

}
