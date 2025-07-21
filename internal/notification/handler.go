package notification

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	http_helper "food-delivery-app-server/pkg/http"
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
	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	notifications, err := h.service.GetUserNotifications(userId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, notifications)
}

func (h *Handler) MarkNotificationAsRead(c *gin.Context) {
	notificationId := c.Param("id")
	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.service.MarkNotificationAsRead(notificationId, userId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Notification marked as read successfully"})
}

func (h *Handler) MarkAllNotificationsAsRead(c *gin.Context) {
	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.service.MarkAllNotificationsAsRead(userId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "All notifications were marked as read"})
}
