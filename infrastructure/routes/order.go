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

	allRoles := order.Group("/")
	{
		allRoles.GET("/:id", orderHandler.GetOrderDetails)
		allRoles.GET("/history", orderHandler.GetOrderHistory)
	}

	ownerAndDriver := order.Group("/", middleware.RequireRoles(models.Owner, models.Driver))
	{
		ownerAndDriver.PUT("/:id", orderHandler.UpdateOrderStatus)
	}

	owner := order.Group("/", middleware.RequireRoles(models.Owner))
	{
		owner.GET("/restaurant/:id", orderHandler.GetOrderByRestaurant)
	}

	customer := order.Group("/", middleware.RequireRoles(models.Customer))
	{
		customer.POST("/restaurant/:id", orderHandler.PlaceOrder)
		customer.PUT("/:id/cancel", orderHandler.CancelOrder)
	}

	driver := order.Group("/", middleware.RequireRoles(models.Driver))
	{
		driver.GET("/available", orderHandler.GetAvailableOrders)
		driver.GET("/assigned", orderHandler.GetAssignedOrders)
	}
}
