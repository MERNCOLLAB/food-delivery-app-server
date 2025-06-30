package infrastructure

import (
	"food-delivery-app-server/middleware"
	"food-delivery-app-server/models"

	"food-delivery-app-server/internal/auth"
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

	restaurantHandler := restaurant.NewHandler(DB)
	ownerGroup := r.Group("/owner", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Owner))
	restaurants := ownerGroup.Group("/restaurant")
	{
		restaurants.POST("/add", restaurantHandler.CreateRestaurant)
		restaurants.GET("/", restaurantHandler.GetRestaurantByOwner)
		restaurants.PUT("/:id", restaurantHandler.UpdateRestaurant)
		restaurants.DELETE("/:id", restaurantHandler.DeleteRestaurant)
	}
}
