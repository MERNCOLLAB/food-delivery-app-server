package user

import (
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
	c.JSON(200, gin.H{
		"message": "Update User Endpoint",
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
