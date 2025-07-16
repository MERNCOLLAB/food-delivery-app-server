package order

import (
	http_helper "food-delivery-app-server/pkg/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAvailableOrders(c *gin.Context) {
	orders, err := h.service.GetAvailableOrders()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"availableOrders": orders})
}

func (h *Handler) GetAssignedOrders(c *gin.Context) {
	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	orders, err := h.service.GetAssignedOrders(userId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"assignedOrders": orders})
}
