package routes

import (
	"food-delivery-app-server/internal/notification"
	"food-delivery-app-server/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterNotificationRoutes(r *gin.Engine, DB *gorm.DB) {
	notifHandler := notification.NewHandler(DB)

	notifGroup := r.Group("/notifications", middleware.JWTAuthMiddleware())
	{
		notifGroup.GET("/", notifHandler.GetUserNotifications)
		notifGroup.PUT("/:id/read", notifHandler.MarkNotificationAsRead)          //not yet functional
		notifGroup.PUT("/mark-all-read", notifHandler.MarkAllNotificationsAsRead) //not yet functional
	}
}
