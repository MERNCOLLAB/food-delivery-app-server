package order

import "github.com/gin-gonic/gin"

func (h *Handler) AcceptOrderByOwner(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Accept Order By Owner Endpoint"})
}

func (h *Handler) GetOrderByRestaurant(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get Order By Restaurant Endpoint"})
}

func (h *Handler) UpdateOrderByOwner(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update Order Status Endpoint"})
}
