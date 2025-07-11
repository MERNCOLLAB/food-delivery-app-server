package order

import (
	http_helper "food-delivery-app-server/pkg/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PlaceOrder(c *gin.Context) {
	restaurantID := c.Param("id")
	req, err := http_helper.BindJSON[PlaceOrderRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	orderRes, err := h.service.PlaceOrder(restaurantID, *req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Order has been placed. Wait for the shop to accept",
		"order":   orderRes,
	})
}

func (h *Handler) CancelOrder(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Cancel Order Endpoint"})
}
