package infrastructure

import (
	"food-delivery-app-server/middleware"
	"food-delivery-app-server/models"

	"food-delivery-app-server/internal/address"
	"food-delivery-app-server/internal/auth"
	"food-delivery-app-server/internal/menuitem"
	"food-delivery-app-server/internal/order"
	"food-delivery-app-server/internal/resetpassword"
	"food-delivery-app-server/internal/restaurant"
	"food-delivery-app-server/internal/user"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	authHandler := auth.NewHandler(DB)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.SignUp)
		authGroup.POST("/signin", authHandler.SignIn)
		authGroup.POST("/oauth-signup/:provider", authHandler.OAuthSignUp)
		authGroup.POST("/oauth-signin/:provider", authHandler.OAuthSignIn)
		authGroup.POST("/send-otp", authHandler.SendOTP)
		authGroup.POST("/verify-otp", authHandler.VerifyOTP)
		authGroup.POST("/signout", authHandler.SignOut)
	}

	userHandler := user.NewHandler(DB)
	userGroup := r.Group("/user", middleware.JWTAuthMiddleware())
	{
		userGroup.PUT("/update", userHandler.UpdateUser)
		userGroup.PUT("/update/profile-picture",
			middleware.UploadImageValidator("image"),
			userHandler.UpdateProfilePicture)
		userGroup.DELETE("/delete", userHandler.DeleteUser)
		userGroup.GET("/", middleware.RequireRoles(models.Admin), userHandler.GetAllUsers)
	}

	resetPasswordHandler := resetpassword.NewHandler(DB)
	resetPasswordGroup := r.Group("/reset-password")
	{
		resetPasswordGroup.POST("/request", resetPasswordHandler.RequestResetPassword)
		resetPasswordGroup.POST("/verify-code", resetPasswordHandler.VerifyResetCode)
		resetPasswordGroup.PUT("/update", resetPasswordHandler.UpdatePassword)
	}

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
		ownerAndCustAddress.GET("/", addressHandler.GetAddress)
		ownerAndCustAddress.PUT("/:id", addressHandler.UpdateAddress)
		ownerAndCustAddress.DELETE("/:id", addressHandler.DeleteAddress)
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
		ownerOrder.GET("/:id", orderHandler.GetOrderByRestaurant)
		ownerOrder.PUT("/:id", orderHandler.UpdateOrderStatus)
	}

	customerGroup := r.Group("/customer", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Customer))
	customerRestaurants := customerGroup.Group("/restaurant")
	{
		customerRestaurants.GET("/")
		customerRestaurants.GET("/:id/menu-items")
	}

	customerMenuItems := customerGroup.Group("/menu-item")
	{
		customerMenuItems.GET("/restaurant/:id")
		customerMenuItems.GET("/:id")
	}

	customerOrder := customerGroup.Group("/order")
	{
		customerOrder.POST("/restaurant/:id")
		customerOrder.GET("/")
		customerOrder.PUT("/:id/cancel")
	}

	customerAddress := customerGroup.Group("/address")
	{
		customerAddress.POST("/")
	}
}
