package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, DB *gorm.DB) {
	RegisterAuthRoutes(r, DB)
	RegisterUserRoutes(r, DB)
	RegisterResetPasswordRoutes(r, DB)
	RegisterRestaurantRoutes(r, DB)
	RegisterMenuItemsRoutes(r, DB)
}
