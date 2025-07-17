package notification

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{service: service}
}

func (h *Handler) GetUserNotifications(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get User Notification Endpoint"})
}

func (h *Handler) MarkNotificationAsRead(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Mark Notification As Read Endpoint"})
}

func (h *Handler) MarkAllNotificationsAsRead(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Mark All Notification As Read Endpoint"})
}
