package infrastructure

import (
	"food-delivery-app-server/middleware"
	"food-delivery-app-server/models"

	"food-delivery-app-server/internal/address"
	"food-delivery-app-server/internal/menuitem"
	"food-delivery-app-server/internal/order"
	"food-delivery-app-server/internal/restaurant"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	addressHandler := address.NewHandler(DB)
	restaurantHandler := restaurant.NewHandler(DB)
	menuItemHandler := menuitem.NewHandler(DB)
	orderHandler := order.NewHandler(DB)

	ownerAndCustMenuItem := r.Group("/menu-item", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Owner, models.Customer))
	{
		ownerAndCustMenuItem.GET("/restaurant/:id", menuItemHandler.GetMenuItemByRestaurant)
	}

	ownerAndCustAddress := r.Group("/address", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Owner, models.Customer))
	{
		ownerAndCustAddress.GET("/", addressHandler.GetAddress)          //not yet functional
		ownerAndCustAddress.PUT("/:id", addressHandler.UpdateAddress)    //not yet functional
		ownerAndCustAddress.DELETE("/:id", addressHandler.DeleteAddress) //not yet functional
	}

	ownerGroup := r.Group("/owner", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Owner))
	ownerRestaurants := ownerGroup.Group("/restaurant")
	{
		ownerRestaurants.POST("/", middleware.UploadImageValidator("image"), restaurantHandler.CreateRestaurant)
		ownerRestaurants.GET("/", restaurantHandler.GetRestaurantByOwner)
		ownerRestaurants.PUT("/:id", middleware.UploadImageValidator("image", true), restaurantHandler.UpdateRestaurant)
		ownerRestaurants.DELETE("/:id", restaurantHandler.DeleteRestaurant)

		ownerRestaurants.POST("/:id/menu-item", middleware.UploadImageValidator("image"), menuItemHandler.CreateMenuItem)
	}

	ownerMenuItems := ownerGroup.Group("/menu-item")
	{
		ownerMenuItems.PUT("/:id", middleware.UploadImageValidator("image", true), menuItemHandler.UpdateMenuItem)
		ownerMenuItems.DELETE("/:id", menuItemHandler.DeleteMenuItem)
	}

	ownerOrder := ownerGroup.Group("/order")
	{
		ownerOrder.GET("/:id", orderHandler.GetOrderByRestaurant) //not yet functional
		ownerOrder.PUT("/:id", orderHandler.UpdateOrderStatus)    //not yet functional
	}

	customerGroup := r.Group("/customer", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Customer))
	customerRestaurants := customerGroup.Group("/restaurant")
	{
		customerRestaurants.GET("/", restaurantHandler.GetAllRestaurants)                      //not yet functional
		customerRestaurants.GET("/:id/menu-items", restaurantHandler.GetMoreRestaurantDetails) //not yet functional
	}

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
