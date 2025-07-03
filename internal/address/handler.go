package address

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

func (h *Handler) CreateAddress(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Create Address Endpoint"})
}

func (h *Handler) GetAddress(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get Address Endpoint"})
}

func (h *Handler) UpdateAddress(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update Address Endpoint"})
}

func (h *Handler) DeleteAddress(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Delete Address Endpoint"})
}
