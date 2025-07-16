package order

import "github.com/gin-gonic/gin"

func (h *Handler) GetOrderByRestaurant(c *gin.Context) {
	restaurantID := c.Param("id")

	orders, err := h.service.GetOrderByRestaurant(restaurantID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"orders": orders})
}
