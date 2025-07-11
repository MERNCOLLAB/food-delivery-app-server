package order

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

// All Roles
func (h *Handler) GetOrderDetails(c *gin.Context) {
	orderId := c.Param("id")

	order, err := h.service.GetOrderDetails(orderId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"orderDetails": order})
}

// Customer & Driver
func (h *Handler) GetOrderHistory(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get Order History Endpoint"})
}

// Owner & Driver
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update Order Endpoint"})
}
