package routes

import (
	"food-delivery-app-server/internal/user"
	"food-delivery-app-server/middleware"
	"food-delivery-app-server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, DB *gorm.DB) {
	userHandler := user.NewHandler(DB)
	userGroup := r.Group("/users", middleware.JWTAuthMiddleware())
	{
		userGroup.PUT("/update", userHandler.UpdateUser)
		userGroup.PUT("/update/profile-picture",
			middleware.UploadImageValidator("image"),
			userHandler.UpdateProfilePicture)
		userGroup.DELETE("/delete", userHandler.DeleteUser)
		userGroup.GET("/", middleware.RequireRoles(models.Admin), userHandler.GetAllUsers)
		userGroup.GET("/driver/:id")   //Not Yet Functional
		userGroup.GET("/customer/:id") //Not Yet Functional
	}

}
