package auth

import (
	"encoding/json"
	"net/http"
	"regexp"

	"food-delivery-app-server/models"
	appErr "food-delivery-app-server/pkg/errors"
	"food-delivery-app-server/pkg/utils"
)

var defaultProfilePic string = "https://res.cloudinary.com/dowkytkyb/image/upload/v1750666850/default_profile_qbzide.png"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SignUp(req SignUpRequest) (*JWTAuthResponse, string, error) {
	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" ||
		req.Name == "" || req.Bio == "" || req.Phone == "" {
		return nil, "", appErr.NewBadRequest("Missing required fields", nil)
	}

	validPhone := regexp.MustCompile(`^\+63[0-9]{10}$`)
	if !validPhone.MatchString(req.Phone) {
		return nil, "", appErr.NewBadRequest("Invalid phone number format. Use +63XXXXXXXXXX", nil)
	}

	if req.Password != req.ConfirmPassword {
		return nil, "", appErr.NewBadRequest("Passwords do not match", nil)
	}

	existingUser, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, "", appErr.NewBadRequest("Failed to verify if the email exists", err)
	}

	if existingUser != nil {
		return nil, "", appErr.NewBadRequest("User with that email already exists", nil)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to hash the password", err)
	}

	userId := utils.GenerateUUID()

	newUser := &models.User{
		ID:             userId,
		Name:           req.Name,
		Email:          req.Email,
		Password:       hashedPassword,
		ProfilePicture: defaultProfilePic,
		Bio:            req.Bio,
		Phone:          req.Phone,
		Role:           models.Role(req.Role),
	}

	createdUser, err := s.repo.CreateUser(newUser)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to create user at database", err)
	}

	token, err := utils.GenerateJWT(createdUser)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to generate token", err)
	}

	userResponse := JWTAuthResponse{
		ID:   createdUser.ID.String(),
		Role: string(createdUser.Role),
	}

	return &userResponse, token, nil
}

func (s *Service) SignIn(req SignInRequest) (*JWTAuthResponse, string, error) {
	if req.Email == "" || req.Password == "" {
		return nil, "", appErr.NewBadRequest("Missing required fields", nil)
	}

	user, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, "", appErr.NewNotFound("Failed to verify if user exists", err)
	}
	if user == nil {
		return nil, "", appErr.NewBadRequest("Invalid email or password", nil)
	}

	if err := utils.ValidatePassword(user.Password, req.Password); err != nil {
		return nil, "", appErr.NewBadRequest("Invalid email or password", err)
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to generate token", err)
	}

	userResponse := JWTAuthResponse{
		ID:   user.ID.String(),
		Role: string(user.Role),
	}

	return &userResponse, token, nil
}

func (s *Service) OAuth(req OAuthRequest, provider string) (*JWTAuthResponse, string, error) {
	var info *UserInfo
	var err error

	switch provider {
	case "google":
		info, err = s.VerifyGoogleToken(req.AccessToken)
	case "facebook":
		info, err = s.VerifyFacebookToken(req.AccessToken)
	default:
		return nil, "", appErr.NewBadRequest("Unsupported provider", nil)
	}

	if err != nil {
		return nil, "", appErr.NewBadRequest("Failed to verify token", err)
	}

	user, err := s.repo.FindUserByEmail(info.Email)
	if err != nil {
		return nil, "", appErr.NewBadRequest("User with that email already exists", nil)
	}

	newUserID := utils.GenerateUUID()

	// Missing Data:  Role, Phone, Number (Email is null for Facebook)
	if user == nil {
		user = &models.User{
			ID:             newUserID,
			Email:          info.Email,
			Name:           info.Name,
			ProfilePicture: info.ProfilePicture,
			Provider:       info.Provider,
		}
		user, err = s.repo.CreateUser(user)
		if err != nil {
			return nil, "", appErr.NewInternal("Failed to create user", err)
		}
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to generate token", err)
	}

	userResponse := JWTAuthResponse{
		ID:   user.ID.String(),
		Role: string(user.Role),
	}

	return &userResponse, token, nil

}

func (s *Service) VerifyGoogleToken(accessToken string) (*UserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + accessToken)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, appErr.NewInternal("failed to verify Google token", err)
	}
	defer resp.Body.Close()

	var data struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &UserInfo{
		Email:          data.Email,
		Name:           data.Name,
		ProfilePicture: data.Picture,
		Provider:       "google",
	}, nil
}

func (s *Service) VerifyFacebookToken(accessToken string) (*UserInfo, error) {
	url := "https://graph.facebook.com/me?fields=id,name,email,picture&access_token=" + accessToken
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, appErr.NewInternal("failed to verify Facebook token", err)
	}
	defer resp.Body.Close()

	var data struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &UserInfo{
		Email:          data.Email,
		Name:           data.Name,
		ProfilePicture: data.Picture.Data.URL,
		Provider:       "facebook",
	}, nil
}
