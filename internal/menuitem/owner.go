package menuitem

import "github.com/gin-gonic/gin"

func (h *Handler) CreateMenuItem(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Create Menu Item Endpoint"})
}

func (h *Handler) UpdateMenuItem(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update Menu Item Endpoint"})
}

func (h *Handler) DeleteMenuItem(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Delete Menu Item Endpoint"})
}
