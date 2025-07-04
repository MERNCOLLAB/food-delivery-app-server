package routes

import (
	"food-delivery-app-server/internal/address"
	"food-delivery-app-server/middleware"
	"food-delivery-app-server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAddressRoutes(r *gin.Engine, DB *gorm.DB) {
	addressHandler := address.NewHandler(DB)

	addresses := r.Group("/addresses", middleware.JWTAuthMiddleware())
	customer := addresses.Group("/", middleware.RequireRoles(models.Customer))
	{
		customer.POST("/", addressHandler.CreateAddress) //not yet functional
	}

	ownerAndCust := addresses.Group("/", middleware.RequireRoles(models.Owner, models.Customer))
	{
		ownerAndCust.GET("/", addressHandler.GetAddress)          //not yet functional
		ownerAndCust.PUT("/:id", addressHandler.UpdateAddress)    //not yet functional
		ownerAndCust.DELETE("/:id", addressHandler.DeleteAddress) //not yet functional
	}
}
