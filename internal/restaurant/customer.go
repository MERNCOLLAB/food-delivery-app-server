package restaurant

import "github.com/gin-gonic/gin"

func (h *Handler) GetAllRestaurants(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get ALl Restaurants Endpoint"})
}

func (h *Handler) GetMoreRestaurantDetails(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get More Restaurant Details Endpoint"})
}
