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
		UserId:      userId,
		Name:        req.Name,
		Description: req.Description,
		Phone:       req.Phone,
		ImageFile:   imageFile.(multipart.File),
		ImageHeader: imageHeader.(*multipart.FileHeader),
	}

	newRestaurant, err := h.service.CreateRestaurant(createRestaurantData)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message":    "Create Restaurant Endpoint",
		"restaurant": newRestaurant,
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
