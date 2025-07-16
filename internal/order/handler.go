package order

import (
	http_helper "food-delivery-app-server/pkg/http"

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
	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	orders, err := h.service.GetOrderHistory(userId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"orderHistory": orders})
}

// Owner & Driver
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	orderId := c.Param("id")

	req, err := http_helper.BindJSON[UpdateOrderStatusRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.UpdateOrderStatus(*req, orderId); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Order status had been updated"})
}
