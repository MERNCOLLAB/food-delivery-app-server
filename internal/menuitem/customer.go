package menuitem

import "github.com/gin-gonic/gin"

func (h *Handler) GetMoreMenuItemDetails(c *gin.Context) {
	menuItemId := c.Param("id")

	menuItem, err := h.service.GetMoreMenuItemDetails(menuItemId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, menuItem)
}
