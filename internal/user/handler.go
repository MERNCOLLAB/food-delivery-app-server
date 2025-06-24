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
		"message": "Update User Endpoint",
		"user": updatedUser,
	})
}

func (h *Handler) UpdateProfilePicture(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Update Profile Picture Endpoint",
	})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Delete User Endpoint",
	})
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get All Users Endpoint",
	})
}
