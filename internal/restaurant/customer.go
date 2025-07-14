package restaurant

import "github.com/gin-gonic/gin"

func (h *Handler) GetAllRestaurants(c *gin.Context) {
	restaurants, err := h.service.GetAllRestaurants()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, restaurants)
}

func (h *Handler) GetMoreRestaurantDetails(c *gin.Context) {
	restaurantId := c.Param("id")

	restaurant, err := h.service.GetMoreRestaurantDetails(restaurantId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"restaurant": restaurant})
}
