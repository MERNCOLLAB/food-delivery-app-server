package routes

import (
	"food-delivery-app-server/internal/menuitem"
	"food-delivery-app-server/middleware"
	"food-delivery-app-server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMenuItemsRoutes(r *gin.Engine, DB *gorm.DB) {
	menuItemHandler := menuitem.NewHandler(DB)

	menuItem := r.Group("/menu-items", middleware.JWTAuthMiddleware())
	owner := menuItem.Group("/", middleware.RequireRoles(models.Owner))
	{
		owner.POST("/restaurant/:id", middleware.UploadImageValidator("image"), menuItemHandler.CreateMenuItem)
		owner.PUT("/:id", middleware.UploadImageValidator("image", true), menuItemHandler.UpdateMenuItem)
		owner.DELETE("/:id", menuItemHandler.DeleteMenuItem)
	}

	customer := menuItem.Group("/", middleware.RequireRoles(models.Customer))
	{
		customer.GET("/:id", menuItemHandler.GetMoreMenuItemDetails)
	}

	ownerAndCust := menuItem.Group("/", middleware.RequireRoles(models.Owner, models.Customer))
	{
		ownerAndCust.GET("/restaurant/:id", menuItemHandler.GetMenuItemByRestaurant)
	}
}
