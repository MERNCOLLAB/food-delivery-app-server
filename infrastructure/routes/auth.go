package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"food-delivery-app-server/internal/auth"
)

func RegisterAuthRoutes(r *gin.Engine, DB *gorm.DB, rdb *redis.Client) {
	authHandler := auth.NewHandler(DB, rdb)

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

}
