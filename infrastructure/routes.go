package infrastructure

import (
	"food-delivery-app-server/internal/order"
	"food-delivery-app-server/middleware"
	"food-delivery-app-server/models"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	orderHandler := order.NewHandler(DB)

	ownerGroup := r.Group("/owner", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Owner))

	ownerOrder := ownerGroup.Group("/order")
	{
		ownerOrder.GET("/:id", orderHandler.GetOrderByRestaurant) //not yet functional
		ownerOrder.PUT("/:id", orderHandler.UpdateOrderStatus)    //not yet functional
	}

	customerGroup := r.Group("/customer", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Customer))

	customerOrder := customerGroup.Group("/order")
	{
		customerOrder.POST("/restaurant/:id", orderHandler.GetOrderByRestaurant) //not yet functional
		customerOrder.GET("/", orderHandler.GetAllPersonalOrders)                //not yet functional
		customerOrder.PUT("/:id/cancel", orderHandler.CancelOrder)               //not yet functional
	}

}
