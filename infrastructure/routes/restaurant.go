package routes

import (
	"food-delivery-app-server/internal/restaurant"
	"food-delivery-app-server/middleware"
	"food-delivery-app-server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRestaurantRoutes(r *gin.Engine, DB *gorm.DB) {
	restaurantHandler := restaurant.NewHandler(DB)

	restaurant := r.Group("/restaurants", middleware.JWTAuthMiddleware())
	owner := restaurant.Group("/", middleware.RequireRoles(models.Owner))
	{
		owner.POST("/", middleware.UploadImageValidator("image"), restaurantHandler.CreateRestaurant)
		owner.GET("/", restaurantHandler.GetRestaurantByOwner)
		owner.PUT("/:id", middleware.UploadImageValidator("image", true), restaurantHandler.UpdateRestaurant)
		owner.DELETE("/:id", restaurantHandler.DeleteRestaurant)
	}

	customer := restaurant.Group("/", middleware.RequireRoles(models.Customer))
	{
		customer.GET("/", restaurantHandler.GetAllRestaurants)                      //not yet functional
		customer.GET("/:id/menu-items", restaurantHandler.GetMoreRestaurantDetails) //not yet functional
	}
}
