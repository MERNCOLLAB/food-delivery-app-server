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

type GoogleResponseData struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

type FacebookResponseData struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

type UserInfo struct {
	Email          string
	FirstName      string
	LastName       string
	ProfilePicture string
	Provider       string
}
