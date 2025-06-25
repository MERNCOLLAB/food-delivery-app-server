package user

import (
	http_helper "food-delivery-app-server/pkg/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{service: service}
}

func (h *Handler) UpdateUser(c *gin.Context) {
	req, err := http_helper.BindJSON[UpdateUserRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	updatedUser, err := h.service.UpdateUser(*req, userId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message": "User has been updated successfully",
		"user":    updatedUser,
	})
}

func (h *Handler) UpdateProfilePicture(c *gin.Context) {
	imageFile, fileExists := c.Get("imageFile")
	imageHeader, headerExists := c.Get("imageHeader")

	if !fileExists || !headerExists {
		c.JSON(400, gin.H{"error": "Image not found in the context"})
		return
	}

	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	updateProfilePicData := UpdateProfilePictureRequest{
		userId:      userId,
		imageFile:   imageFile,
		imageHeader: imageHeader,
	}

	url, err := h.service.UpdateProfilePicture(updateProfilePicData)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Profile picture has been updated successfully",
		"url":     url,
	})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(400, gin.H{"error": "email query parameter is required"})
		return
	}

	userId, err := http_helper.ExtractUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.service.DeleteUser(userId, email)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message": "User account has been deleted successfully",
	})
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"allUsers": users,
	})
}
