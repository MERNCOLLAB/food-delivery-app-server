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
		Address:     req.Address,
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
	restaurantId := c.Param("id")
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

	req, err := http_helper.BindFormJSON[UpdateRestaurantRequest](c, "data")
	if err != nil {
		c.Error(err)
		return
	}

	updateRestaurantData := UpdateRestaurantRequest{
		Name:        req.Name,
		Description: req.Description,
		Phone:       req.Phone,
		Address:     req.Address,
		ImageFile:   imgFilePtr,
		ImageHeader: imgHeaderPtr,
	}

	updatedRestaurant, err := h.service.UpdateRestaurant(restaurantId, updateRestaurantData)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message":    "Restaurant has been updated successfully",
		"restaurant": updatedRestaurant,
	})
}

func (h *Handler) DeleteRestaurant(c *gin.Context) {
	restauarantId := c.Param("id")

	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.service.DeleteRestaurant(userId, restauarantId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Restaurant has been deleted successfully",
	})
}
