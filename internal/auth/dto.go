package auth

import "food-delivery-app-server/models"

// JWT Authentication Feature
type Role string

const (
	Admin    Role = "ADMIN"
	Owner    Role = "OWNER"
	Driver   Role = "DRIVER"
	Customer Role = "CUSTOMER"
)

type SignUpRequest struct {
	FirstName       string
	LastName        string
	Email           string
	Bio             string
	Phone           string
	Address         string
	Password        string
	ConfirmPassword string
	Token           string
	Role            Role
}

type JWTAuthResponse struct {
	ID   string
	Role string
}

type SignInRequest struct {
	Email    string
	Password string
}

// Oauth Feature
type OAuthRequest struct {
	AccessToken string `json:"accessToken"`
}

// Send and Validate Phone OTP
type SendOTPRequest struct {
	Phone string `json:"phone"`
}

type VerifyOTPRequest struct {
	RedisKey string `json:"redisKey"`
	Phone    string `json:"phone"`
	OTP      string `json:"otp"`
}

// Admin Actions
type SendSignUpFormRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}

type SignUpDecisionRequest struct {
	IsAccepted *bool `json:"isAccepted" binding:"required"`
}

type PendingSignUp struct {
	User    *models.User    `json:"user"`
	Address *models.Address `json:"address"`
}
