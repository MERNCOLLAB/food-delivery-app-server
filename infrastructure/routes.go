package infrastructure

import (
	"food-delivery-app-server/middleware"
	"food-delivery-app-server/models"

	"food-delivery-app-server/internal/address"
	"food-delivery-app-server/internal/menuitem"
	"food-delivery-app-server/internal/order"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	addressHandler := address.NewHandler(DB)

	menuItemHandler := menuitem.NewHandler(DB)
	orderHandler := order.NewHandler(DB)

	ownerAndCustAddress := r.Group("/address", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Owner, models.Customer))
	{
		ownerAndCustAddress.GET("/", addressHandler.GetAddress)          //not yet functional
		ownerAndCustAddress.PUT("/:id", addressHandler.UpdateAddress)    //not yet functional
		ownerAndCustAddress.DELETE("/:id", addressHandler.DeleteAddress) //not yet functional
	}

	ownerGroup := r.Group("/owner", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Owner))

	ownerOrder := ownerGroup.Group("/order")
	{
		ownerOrder.GET("/:id", orderHandler.GetOrderByRestaurant) //not yet functional
		ownerOrder.PUT("/:id", orderHandler.UpdateOrderStatus)    //not yet functional
	}

	customerGroup := r.Group("/customer", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Customer))

	customerMenuItems := customerGroup.Group("/menu-item")
	{
		customerMenuItems.GET("/:id", menuItemHandler.GetMoreMenuItemDetails) //not yet functional
	}

	customerOrder := customerGroup.Group("/order")
	{
		customerOrder.POST("/restaurant/:id", orderHandler.GetOrderByRestaurant) //not yet functional
		customerOrder.GET("/", orderHandler.GetAllPersonalOrders)                //not yet functional
		customerOrder.PUT("/:id/cancel", orderHandler.CancelOrder)               //not yet functional
	}

	customerAddress := customerGroup.Group("/address")
	{
		customerAddress.POST("/", addressHandler.CreateAddress) //not yet functional
	}
}
