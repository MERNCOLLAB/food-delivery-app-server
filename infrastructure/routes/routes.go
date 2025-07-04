package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, DB *gorm.DB) {
	// Baseline Feature Routes
	RegisterAuthRoutes(r, DB)
	RegisterUserRoutes(r, DB)
	RegisterResetPasswordRoutes(r, DB)

	// Core Routes
	RegisterRestaurantRoutes(r, DB)
	RegisterMenuItemsRoutes(r, DB)
	RegisterAddressRoutes(r, DB)
	RegisterOrderRoutes(r, DB)
}
