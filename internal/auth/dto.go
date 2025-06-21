package auth

type Role string

const (
	Admin    Role = "ADMIN"
	Owner    Role = "OWNER"
	Driver   Role = "DRIVER"
	Customer Role = "CUSTOMER"
)

type SignUpRequest struct {
	Name            string
	Email           string
	Bio             string
	Phone           string
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
