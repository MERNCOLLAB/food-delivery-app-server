package order

import "github.com/gin-gonic/gin"

func (h *Handler) GetAvailableOrders(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get Available Orders Endpoint"})
}

func (h *Handler) GetAssignedOrders(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get Assigned Orders Endpoint"})
}

func (h *Handler) UpdateOrderByDriver(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update Order Status"})
}
