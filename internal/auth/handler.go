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
		"message":  "Your application has been sent to admin for approval",
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

	redisKey, err := h.service.OAuthSignUp(*req, provider)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"redisKey": redisKey,
		"message":  fmt.Sprintf("OAuth with %s succeeded. Please proceed with providing your phone number", provider),
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

func (h *Handler) SendOTPToPhone(c *gin.Context) {
	redisKey := c.Param("id")
	req, err := http_helper.BindJSON[SendOTPRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.SendOTPToPhone(redisKey, req.Phone); err != nil {
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

func (h *Handler) SignUpDecision(c *gin.Context) {
	signUpID := c.Param("id")
	req, err := http_helper.BindJSON[SignUpDecisionRequest](c)
	if err != nil {
		c.Error(err)
		return
	}

	userResData, token, err := h.service.SignUpDecision(*req, signUpID)
	if err != nil {
		c.Error(err)
		return
	}

	if req.IsAccepted == nil {
		c.JSON(400, gin.H{"error": "isAccepted is required"})
		return
	}

	if !*req.IsAccepted {
		c.JSON(200, gin.H{"message": "Sign Up Application has been rejected"})
		return
	}

	utils.SetCookie(c, token, 3600*5)

	c.JSON(200, gin.H{
		"message": "Sign Up Application has been decided",
		"user":    userResData,
	})
}
