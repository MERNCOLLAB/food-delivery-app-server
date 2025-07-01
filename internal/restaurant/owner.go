package restaurant

import (
	http_helper "food-delivery-app-server/pkg/http"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateRestaurant(c *gin.Context) {
	imageFile, fileExists := c.Get("imageFile")
	imageHeader, headerExists := c.Get("imageHeader")

	if !fileExists || !headerExists {
		c.JSON(400, gin.H{"error": "Image not found in the context"})
		return
	}

	req, err := http_helper.BindFormJSON[CreateRestaurantRequest](c, "data")
	if err != nil {
		c.Error(err)
		return
	}

	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	createRestaurantData := CreateRestaurantRequest{
		Name:        req.Name,
		Description: req.Description,
		Phone:       req.Phone,
		ImageFile:   imageFile.(multipart.File),
		ImageHeader: imageHeader.(*multipart.FileHeader),
	}

	newRestaurant, err := h.service.CreateRestaurant(userId, createRestaurantData)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message":    "Restaurant has been successfully added",
		"restaurant": newRestaurant,
	})
}

func (h *Handler) GetRestaurantByOwner(c *gin.Context) {
	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	restaurantList, err := h.service.GetRestaurantByOwner(userId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"restaurants": restaurantList,
	})
}

func (h *Handler) UpdateRestaurant(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "Restaurant has been updated successfully",
	})
}

func (h *Handler) DeleteRestaurant(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Delete Restaurant Endpoint",
	})
}
