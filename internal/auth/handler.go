package auth

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

func (h *Handler) SignUp(c *gin.Context) {

	err := h.service.TestError()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Sign Up Endpoint",
	})
}

func (h *Handler) SignIn(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Sign In Endpoint",
	})
}

func (h *Handler) OAuth(c *gin.Context) {

}

func (h *Handler) SignOut(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Sign Out Endpoint",
	})
}
