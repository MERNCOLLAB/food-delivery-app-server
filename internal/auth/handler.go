package auth

import (
	http_helper "food-delivery-app-server/pkg/http"
	"food-delivery-app-server/pkg/utils"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(db *gorm.DB, rdb *redis.Client) *Handler {
	repo := NewRepository(db)
	service := NewService(repo, rdb)
	return &Handler{service: service}
}

func (h *Handler) SignUp(c *gin.Context) {
	req, err := http_helper.BindJSON[SignUpRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	signUpID, err := h.service.SignUp(*req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message":  "You have successfully registered an account",
		"signUpID": signUpID,
	})
}

func (h *Handler) SignIn(c *gin.Context) {
	req, err := http_helper.BindJSON[SignInRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	user, token, err := h.service.SignIn(*req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SetCookie(c, token, 3600*5)

	c.JSON(200, gin.H{
		"message": "Signed In Successfully",
		"user":    user,
	})
}

func (h *Handler) OAuthSignUp(c *gin.Context) {
	provider := c.Param("provider")
	req, err := http_helper.BindJSON[OAuthRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	stateID, err := h.service.OAuthSignUp(*req, provider)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"stateId": stateID,
		"message": fmt.Sprintf("OAuth with %s succeeded; please verify phone next.", provider),
	})
}

func (h *Handler) OAuthSignIn(c *gin.Context) {
	provider := c.Param("provider")
	req, err := http_helper.BindJSON[OAuthRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	user, token, err := h.service.OAuthSignIn(*req, provider)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SetCookie(c, token, 3600*5)

	c.JSON(200, gin.H{
		"message": "You have successfully signed in",
		"user":    user,
	})
}

func (h *Handler) SendOTP(c *gin.Context) {
	req, err := http_helper.BindJSON[SendOTPRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.SendOTP(req.StateID, req.Phone); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "OTP sent to your phone number"})
}

func (h *Handler) VerifyOTP(c *gin.Context) {
	req, err := http_helper.BindJSON[VerifyOTPRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	user, token, err := h.service.VerifyOTP(*req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SetCookie(c, token, 3600*5)

	c.JSON(200, gin.H{
		"message": "OTP verified and user created successfully",
		"user":    user,
	})
}

func (h *Handler) SignOut(c *gin.Context) {
	utils.ClearCookie(c)

	c.JSON(200, gin.H{
		"message": "You have signed out successfully",
	})
}

func (h *Handler) SendSignUpForm(c *gin.Context) {
	req, err := http_helper.BindJSON[SendSignUpFormRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.SendSignUpForm(*req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"message": "A sign up form invitation with link was sent to the provided email",
	})
}
