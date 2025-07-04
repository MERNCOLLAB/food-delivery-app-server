package routes

import (
	"food-delivery-app-server/internal/order"
	"food-delivery-app-server/middleware"
	"food-delivery-app-server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterOrderRoutes(r *gin.Engine, DB *gorm.DB) {
	orderHandler := order.NewHandler(DB)

	order := r.Group("/orders", middleware.JWTAuthMiddleware())
	owner := order.Group("/", middleware.RequireRoles(models.Owner))
	{
		owner.GET("/restaurant/:id", orderHandler.GetOrderByRestaurant) //not yet functional
		owner.PUT("/:id", orderHandler.UpdateOrderStatus)               //not yet functional
	}

	customer := order.Group("/", middleware.RequireRoles(models.Customer))
	{
		customer.POST("/restaurant/:id", orderHandler.PlaceOrder)
		customer.GET("/", orderHandler.GetAllPersonalOrders)
		customer.PUT("/cancel/:id", orderHandler.CancelOrder)
	}
}
