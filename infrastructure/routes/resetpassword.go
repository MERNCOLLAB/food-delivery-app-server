package routes

import (
	"food-delivery-app-server/internal/resetpassword"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterResetPasswordRoutes(r *gin.Engine, DB *gorm.DB) {
	resetPasswordHandler := resetpassword.NewHandler(DB)
	resetPasswordGroup := r.Group("/reset-password")
	{
		resetPasswordGroup.POST("/request", resetPasswordHandler.RequestResetPassword)
		resetPasswordGroup.POST("/verify-code", resetPasswordHandler.VerifyResetCode)
		resetPasswordGroup.PUT("/update", resetPasswordHandler.UpdatePassword)
	}

}
