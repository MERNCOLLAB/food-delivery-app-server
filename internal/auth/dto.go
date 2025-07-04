package auth

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
	StateID string `json:"stateId"`
	Phone   string `json:"phone"`
}

type VerifyOTPRequest struct {
	StateID string `json:"stateId"`
	Phone   string `json:"phone"`
	OTP     string `json:"otp"`
}
