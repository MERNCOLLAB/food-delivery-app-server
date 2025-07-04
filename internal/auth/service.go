package auth

import (
	"context"
	"strings"
	"time"

	"food-delivery-app-server/models"
	appErr "food-delivery-app-server/pkg/errors"
	"food-delivery-app-server/pkg/geocode"
	"food-delivery-app-server/pkg/oauth"
	"food-delivery-app-server/pkg/sms"
	"food-delivery-app-server/pkg/utils"
)

var DefaultProfilePic string = "https://res.cloudinary.com/dowkytkyb/image/upload/v1750666850/default_profile_qbzide.png"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SignUp(req SignUpRequest) (*JWTAuthResponse, string, error) {
	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" || req.Address == "" ||
		req.FirstName == "" || req.LastName == "" || req.Bio == "" || req.Phone == "" {
		return nil, "", appErr.NewBadRequest("Missing required fields", nil)
	}

	if err := sms.ValidatePhone(req.Phone); err != nil {
		return nil, "", appErr.NewBadRequest("Invalid Phone Number Format", err)
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
	addressId := utils.GenerateUUID()

	newUser := &models.User{
		ID:             userId,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Password:       hashedPassword,
		ProfilePicture: DefaultProfilePic,
		Bio:            req.Bio,
		Phone:          req.Phone,
		Role:           models.Role(req.Role),
	}

	ctx := context.Background()
	lat, long, err := geocode.Geocode(ctx, req.Address)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to geocode the provided address", err)
	}

	newAddress := &models.Address{
		ID:        addressId,
		UserID:    &userId,
		Address:   req.Address,
		IsDefault: true,
		Latitude:  lat,
		Longitude: long,
	}

	createdUser, err := s.repo.CreateUser(newUser)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to create user at database", err)
	}

	_, err = s.repo.CreateAddress(newAddress)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to create address at database", err)
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

func (s *Service) OAuthSignUp(req OAuthRequest, provider string) (string, error) {
	var info *oauth.UserInfo
	var err error

	switch provider {
	case "google":
		info, err = oauth.VerifyGoogleToken(req.AccessToken)
	case "facebook":
		info, err = oauth.VerifyFacebookToken(req.AccessToken)
	default:
		return "", appErr.NewBadRequest("Unsupported provider", nil)
	}

	if err != nil {
		return "", appErr.NewBadRequest("Failed to verify token", err)
	}

	stateID := utils.GenerateStateID(info)

	return stateID, nil
}

func (s *Service) OAuthSignIn(req OAuthRequest, provider string) (*JWTAuthResponse, string, error) {
	var info *oauth.UserInfo
	var user *models.User
	var err error

	// For retrieving user data (info) from OAuth provder
	switch provider {
	case "google":
		info, err = oauth.VerifyGoogleToken(req.AccessToken)
	case "facebook":
		info, err = oauth.VerifyFacebookToken(req.AccessToken)
	default:
		return nil, "", appErr.NewBadRequest("Unsupported provider", nil)
	}
	if err != nil {
		return nil, "", appErr.NewBadRequest("Failed to verify OAuth token", err)
	}

	// For validating if user account exists in the database
	switch provider {
	case "google":
		user, err = s.repo.FindUserByEmail(info.Email)
	case "facebook":
		if strings.HasPrefix(info.ProfilePicture, "https://platform-lookaside.fbsbx.com") {
			user, err = s.repo.FindFacebookUserByProfilePicturePrefix(info.ProfilePicture)
		} else {
			return nil, "", appErr.NewBadRequest("Invalid Facebook profile picture", nil)
		}
	default:
		return nil, "", appErr.NewBadRequest("Unsupported provider", nil)
	}

	if err != nil {
		return nil, "", appErr.NewInternal("Account not found, sign up first", err)
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

func (s *Service) SendOTP(stateID, phone string) error {
	if err := sms.ValidatePhone(phone); err != nil {
		return appErr.NewBadRequest("Invalid phone number", err)
	}

	utils.OAuthTempStore.RLock()
	data, ok := utils.OAuthTempStore.M[stateID]
	utils.OAuthTempStore.RUnlock()
	if !ok || time.Now().After(data.ExpiresAt) {
		return appErr.NewBadRequest("Invalid or Expired State ID", nil)
	}

	otp := utils.GenerateOTP()

	utils.OtpStore.Lock()
	utils.OtpStore.M[phone] = otp
	utils.OtpStore.Unlock()

	if err := sms.SendOTPTextBee(phone, otp); err != nil {
		return appErr.NewInternal("Failed to send OTP via SMS", err)
	}

	return nil
}

func (s *Service) VerifyOTP(req VerifyOTPRequest) (*JWTAuthResponse, string, error) {
	phone := req.Phone
	otp := req.OTP
	stateID := req.StateID

	if err := sms.ValidatePhone(phone); err != nil {
		return nil, "", appErr.NewBadRequest("Invalid phone number", err)
	}

	utils.OtpStore.RLock()
	expectedOTP, ok := utils.OtpStore.M[phone]
	utils.OtpStore.RUnlock()
	if !ok || expectedOTP != otp {
		return nil, "", appErr.NewBadRequest("Invalid or expired OTP", nil)
	}

	utils.OAuthTempStore.RLock()
	oAuthData, ok := utils.OAuthTempStore.M[stateID]
	utils.OAuthTempStore.RUnlock()
	if !ok || time.Now().After(oAuthData.ExpiresAt) {
		return nil, "", appErr.NewBadRequest("Invalid or expired state ID", nil)
	}

	info, ok := oAuthData.Info.(*oauth.UserInfo)
	if !ok {
		return nil, "", appErr.NewInternal("Failed to parse OAuth user info", nil)
	}

	userId := utils.GenerateUUID()
	newUser := &models.User{
		ID:             userId,
		FirstName:      info.FirstName,
		LastName:       info.LastName,
		Email:          info.Email,
		ProfilePicture: info.ProfilePicture,
		Bio:            "",
		Phone:          phone,
		Role:           models.Customer,
		Provider:       info.Provider,
	}

	createdUser, err := s.repo.CreateUser(newUser)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to create user at database", err)
	}

	token, err := utils.GenerateJWT(createdUser)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to generate token", err)
	}

	utils.CleanMemory(phone, stateID)

	userResponse := JWTAuthResponse{
		ID:   createdUser.ID.String(),
		Role: string(createdUser.Role),
	}

	return &userResponse, token, nil
}
