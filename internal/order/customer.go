package order

import "github.com/gin-gonic/gin"

func (h *Handler) PlaceOrder(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Place Order Endpoint"})
}

func (h *Handler) GetAllPersonalOrders(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get All Personal Orders Endpoint"})
}

func (h *Handler) CancelOrder(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Cancel Order Endpoint"})
}
