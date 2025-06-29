package resetpassword

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

func (h *Handler) RequestResetPassword(c *gin.Context) {
	req, err := http_helper.BindJSON[ResetPasswordRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	if err = h.service.RequestResetPassword(*req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Reset code has been sent to your email"})
}

func (h *Handler) VerifyResetCode(c *gin.Context) {
	req, err := http_helper.BindJSON[VerifyCodeRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	if err = h.service.VerifyResetCode(*req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Reset Code is valid. Proceed to change your password"})
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	req, err := http_helper.BindJSON[UpdatePasswordRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	email, err := h.service.UpdatePassword(*req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Your password has been updated successfuly", "user": email})
}
