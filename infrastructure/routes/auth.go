package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"food-delivery-app-server/internal/auth"
)

func RegisterAuthRoutes(r *gin.Engine, DB *gorm.DB) {
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

}
