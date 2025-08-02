package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, DB *gorm.DB, rdb *redis.Client) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "🍱 Food Delivery App Server is running. Server developed by Bryan Castillo."})
	})

	// Baseline Feature Routes
	RegisterAuthRoutes(r, DB, rdb)
	RegisterUserRoutes(r, DB)
	RegisterResetPasswordRoutes(r, DB)

	// Core Routes
	RegisterRestaurantRoutes(r, DB)
	RegisterMenuItemsRoutes(r, DB)
	RegisterAddressRoutes(r, DB)
	RegisterOrderRoutes(r, DB)
	RegisterNotificationRoutes(r, DB)
}
