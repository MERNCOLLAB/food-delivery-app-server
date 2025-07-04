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
		Category:    req.Category,
		ImageFile:   imageFile.(multipart.File),
		ImageHeader: imageHeader.(*multipart.FileHeader),
	}

	newMenuItem, err := h.service.CreateMenuItem(restaurantId, createMenuItemData)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Your food item has been added successfully", "menuItem": newMenuItem})
}

func (h *Handler) UpdateMenuItem(c *gin.Context) {
	menuItemId := c.Param("id")
	imageFile, fileExists := c.Get("imageFile")
	imageHeader, headerExists := c.Get("imageHeader")

	var imgFilePtr *multipart.File
	var imgHeaderPtr *multipart.FileHeader

	if fileExists && headerExists {
		f := imageFile.(multipart.File)
		h := imageHeader.(*multipart.FileHeader)

		imgFilePtr = &f
		imgHeaderPtr = h
	} else {
		imgFilePtr = nil
		imgHeaderPtr = nil
	}

	req, err := http_helper.BindFormJSON[UpdateMenuItemRequest](c, "data")
	if err != nil {
		c.Error(err)
		return
	}

	updateMenuItemData := UpdateMenuItemRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		ImageFile:   imgFilePtr,
		ImageHeader: imgHeaderPtr,
	}

	updatedMenuItem, err := h.service.UpdateMenuItem(menuItemId, updateMenuItemData)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message":  "Food item has been updated successfully",
		"menuItem": updatedMenuItem,
	})
}

func (h *Handler) DeleteMenuItem(c *gin.Context) {
	menuItemId := c.Param("id")

	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.service.DeleteMenuItem(userId, menuItemId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Food item has been deleted successfully"})
}
