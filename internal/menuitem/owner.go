package menuitem

import (
	http_helper "food-delivery-app-server/pkg/http"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateMenuItem(c *gin.Context) {
	imageFile, fileExists := c.Get("imageFile")
	imageHeader, headerExists := c.Get("imageHeader")

	if !fileExists || !headerExists {
		c.JSON(400, gin.H{"error": "Image not found in the context"})
		return
	}

	restaurantId := c.Param("id")

	req, err := http_helper.BindFormJSON[CreateMenuItemRequest](c, "data")
	if err != nil {
		c.Error(err)
		return
	}

	createMenuItemData := CreateMenuItemRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageFile:   imageFile.(multipart.File),
		ImageHeader: imageHeader.(*multipart.FileHeader),
	}

	newMenuItem, err := h.service.CreateMenuItem(restaurantId, createMenuItemData)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Create Menu Item Endpoint", "menuItem": newMenuItem})
}

func (h *Handler) UpdateMenuItem(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update Menu Item Endpoint"})
}

func (h *Handler) DeleteMenuItem(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Delete Menu Item Endpoint"})
}
