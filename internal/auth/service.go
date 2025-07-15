package auth

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"food-delivery-app-server/models"
	appErr "food-delivery-app-server/pkg/errors"
	"food-delivery-app-server/pkg/geocode"
	"food-delivery-app-server/pkg/oauth"
	"food-delivery-app-server/pkg/sms"
	"food-delivery-app-server/pkg/utils"

	"github.com/redis/go-redis/v9"
)

var DefaultProfilePic string = "https://res.cloudinary.com/dowkytkyb/image/upload/v1750666850/default_profile_qbzide.png"

type Service struct {
	repo *Repository
	rdb  *redis.Client
}

func NewService(repo *Repository, rdb *redis.Client) *Service {
	return &Service{repo: repo, rdb: rdb}
}

func (s *Service) SignUpDriver(req SignUpRequest) (string, error) {
	// Missing Required Validation
	if req.Email == "" || req.Address == "" ||
		req.FirstName == "" || req.LastName == "" || req.Bio == "" || req.Phone == "" {
		return "", appErr.NewBadRequest("Missing required fields", nil)
	}

	// Validate Phone Format
	if err := sms.ValidatePhone(req.Phone); err != nil {
		return "", appErr.NewBadRequest("Invalid Phone Number Format", err)
	}

	// Existing User Validation
	existingUser, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return "", appErr.NewBadRequest("Failed to verify if the email exists", err)
	}

	if existingUser != nil {
		return "", appErr.NewBadRequest("User with that email already exists", nil)
	}

	// User and Address Data Preparation
	userId := utils.GenerateUUID()
	addressId := utils.GenerateUUID()

	newUser := &models.User{
		ID:             userId,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		ProfilePicture: DefaultProfilePic,
		Bio:            req.Bio,
		Phone:          req.Phone,
		Role:           models.Role(req.Role),
	}

	ctx := context.Background()
	lat, long, err := geocode.Geocode(ctx, req.Address)
	if err != nil {
		return "", appErr.NewInternal("Failed to geocode the provided address", err)
	}

	newAddress := &models.Address{
		ID:        addressId,
		UserID:    &userId,
		Address:   req.Address,
		IsDefault: true,
		Latitude:  lat,
		Longitude: long,
	}

	// Marshal the Data for Redis Storage
	pendingSignUpID := utils.GenerateUUIDStr()
	pendingData := map[string]interface{}{
		"user":    newUser,
		"address": newAddress,
	}

	data, err := json.Marshal(pendingData)
	if err != nil {
		return "", appErr.NewInternal("Failed to serialize pending signup data", err)
	}

	expiry := 30 * time.Minute
	err = s.rdb.Set(ctx, "pending_signup:"+pendingSignUpID, data, expiry).Err()
	if err != nil {
		return "", appErr.NewInternal("Failed to store pending signup in Redis", err)
	}

	return pendingSignUpID, nil
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

	redisKey := utils.SetTempCustomer(s.rdb, info)

	return redisKey, nil
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

func (s *Service) SendOTPToPhone(redisKey, phone string) error {
	if err := sms.ValidatePhone(phone); err != nil {
		return appErr.NewBadRequest("Invalid phone number", err)
	}

	_, err := utils.GetTempUser(s.rdb, redisKey)
	if err != nil {
		return appErr.NewBadRequest("Invalid or expired temporary user data", err)
	}

	otp := utils.GenerateOTP()

	if err := utils.SetOTP(s.rdb, phone, otp, 5*time.Minute); err != nil {
		return appErr.NewInternal("Failed to store OTP", err)
	}

	if err := sms.SendOTPTextBee(phone, otp); err != nil {
		return appErr.NewInternal("Failed to send OTP via SMS", err)
	}

	return nil
}

func (s *Service) VerifyOTP(req VerifyOTPRequest) (*JWTAuthResponse, string, error) {
	phone := req.Phone
	otp := req.OTP
	redisKey := req.RedisKey

	if err := sms.ValidatePhone(phone); err != nil {
		return nil, "", appErr.NewBadRequest("Invalid phone number", err)
	}

	storedOTP, err := utils.GetOTP(s.rdb, phone)
	if err != nil || storedOTP != otp {
		return nil, "", appErr.NewBadRequest("Invalid or expired OTP", nil)
	}

	oAuthData, err := utils.GetTempUser(s.rdb, redisKey)
	if err != nil {
		return nil, "", appErr.NewBadRequest("Invalid or expired redis key", nil)
	}

	info, ok := oAuthData.Info.(*oauth.UserInfo)
	if !ok {
		b, _ := json.Marshal(oAuthData.Info)
		var userInfo oauth.UserInfo
		if err := json.Unmarshal(b, &userInfo); err != nil {
			return nil, "", appErr.NewInternal("Failed to parse OAuth user info", err)
		}
		info = &userInfo
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

	userResponse := JWTAuthResponse{
		ID:   createdUser.ID.String(),
		Role: string(createdUser.Role),
	}

	return &userResponse, token, nil
}
