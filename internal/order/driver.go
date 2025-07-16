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

func (h *Handler) UpdateDriverOrderStatus(c *gin.Context) {
	orderId := c.Param("id")

	driverId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	req, err := http_helper.BindJSON[UpdateOrderStatusRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.UpdateDriverOrderStatus(*req, orderId, driverId); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Driver had updated the order status"})
}
