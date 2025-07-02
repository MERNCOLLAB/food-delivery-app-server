package menuitem

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

func (h *Handler) GetMenuItemByRestaurant(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get Menu Items By Restaurant Endpoint"})
}
