package restaurant

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateRestaurant(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Create Restaurant Endpoint",
	})
}

func (h *Handler) GetRestaurantByOwner(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get Restaurant By Owner Endpoint",
	})
}

func (h *Handler) UpdateRestaurant(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Update Restaurant Endpoint",
	})
}

func (h *Handler) DeleteRestaurant(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Delete Restaurant Endpoint",
	})
}
