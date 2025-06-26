package resetpassword

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

func (h *Handler) RequestResetPassword(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Request Reset Password Endpoint"})
}

func (h *Handler) VerifyResetCode(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Verify Reset Code Endpoint"})
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update Password Endpoint"})
}
