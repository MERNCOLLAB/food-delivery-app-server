package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"food-delivery-app-server/internal/auth"
	"food-delivery-app-server/middleware"
	"food-delivery-app-server/models"
)

func RegisterAuthRoutes(r *gin.Engine, DB *gorm.DB, rdb *redis.Client) {
	authHandler := auth.NewHandler(DB, rdb)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.SignUp)
		authGroup.POST("/signin", authHandler.SignIn)
		authGroup.POST("/oauth-signup/:provider", authHandler.OAuthSignUp)
		authGroup.POST("/oauth-signin/:provider", authHandler.OAuthSignIn)
		authGroup.POST("/send-otp/:id", authHandler.SendOTPToPhone)
		authGroup.POST("/verify-otp", authHandler.VerifyOTP)
		authGroup.POST("/signout", authHandler.SignOut)
	}

	adminOnly := authGroup.Group("/", middleware.JWTAuthMiddleware(), middleware.RequireRoles(models.Admin))
	{
		adminOnly.POST("/send-signup", authHandler.SendSignUpForm)
		adminOnly.POST("/pending-signups/:id/decision", authHandler.SignUpDecision)
	}
}
