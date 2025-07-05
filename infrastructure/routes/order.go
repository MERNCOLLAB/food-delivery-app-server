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
		allRoles.GET("/:id") //not yet functional
	}

	custAndDriver := order.Group("/", middleware.RequireRoles(models.Customer, models.Driver))
	{
		custAndDriver.GET("/history") //not yet functional
	}

	owner := order.Group("/", middleware.RequireRoles(models.Owner))
	{
		owner.GET("/restaurant/:id", orderHandler.GetOrderByRestaurant) //not yet functional
		owner.PUT("/:id", orderHandler.UpdateOrderStatus)               //not yet functional
	}

	customer := order.Group("/", middleware.RequireRoles(models.Customer))
	{
		customer.POST("/restaurant/:id", orderHandler.PlaceOrder) //not yet functional
		customer.GET("/", orderHandler.GetAllPersonalOrders)      //not yet functional
		customer.PUT("/:id/cancel", orderHandler.CancelOrder)     //not yet functional
	}

	driver := order.Group("/", middleware.RequireRoles(models.Driver))
	{
		driver.GET("/available")         //not yet functional
		driver.PUT("/:id/accept")        //not yet functional
		driver.GET("/assigned")          //not yet functional
		driver.PUT("/:id/update-status") //not yet functional
	}
}
