package infrastructure

import (
	"food-delivery-app-server/internal/auth"
	"food-delivery-app-server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	authHandler := auth.NewHandler(DB)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.SignUp)
		authGroup.POST("/signin", authHandler.SignIn)
		authGroup.POST("/oauth/:provider", authHandler.OAuth)
		authGroup.POST("/signout", authHandler.SignOut)
	}
	userGroup := r.Group("/user", middleware.JWTAuthMiddleware())
	{
		userGroup.PUT("/update")
		userGroup.PUT("/update/profile-picture")
		userGroup.DELETE("/delete")
		userGroup.GET("/")
	}
}
