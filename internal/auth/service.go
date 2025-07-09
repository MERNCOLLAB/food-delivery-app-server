package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"food-delivery-app-server/models"
	"food-delivery-app-server/pkg/email"
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

func (s *Service) SignUp(req SignUpRequest) (string, error) {
	// Missing Required Validation
	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" || req.Address == "" ||
		req.FirstName == "" || req.LastName == "" || req.Bio == "" || req.Phone == "" {
		return "", appErr.NewBadRequest("Missing required fields", nil)
	}

	// Validate Phone Format
	if err := sms.ValidatePhone(req.Phone); err != nil {
		return "", appErr.NewBadRequest("Invalid Phone Number Format", err)
	}

	// Password Mismatch
	if req.Password != req.ConfirmPassword {
		return "", appErr.NewBadRequest("Passwords do not match", nil)
	}

	// Existing User Validation
	existingUser, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return "", appErr.NewBadRequest("Failed to verify if the email exists", err)
	}

	if existingUser != nil {
		return "", appErr.NewBadRequest("User with that email already exists", nil)
	}

	// Sign Up Link Token Validation
	if req.Token == "" {
		return "", appErr.NewBadRequest("Missing invitation token", nil)
	}
	val, err := s.rdb.Get(context.Background(), "signup_invite:"+req.Token).Result()

	if err == redis.Nil {
		return "", appErr.NewBadRequest("Invalid or expired invitation token", nil)
	} else if err != nil {
		return "", appErr.NewInternal("Failed to verify invitation token", err)
	}

	var invite SendSignUpFormRequest
	if err := json.Unmarshal([]byte(val), &invite); err != nil {
		return "", appErr.NewInternal("Failed to parse invitation data", err)
	}

	if invite.Email != req.Email {
		return "", appErr.NewBadRequest("Invitation token does not match email", nil)
	}

	// Password Hashing
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", appErr.NewInternal("Failed to hash the password", err)
	}

	// User and Address Data Preparation
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

	expiry := 10 * time.Minute
	err = s.rdb.Set(ctx, "pending_signup:"+pendingSignUpID, data, expiry).Err()
	if err != nil {
		return "", appErr.NewInternal("Failed to store pending signup in Redis", err)
	}

	// Sending Admin Users Notification of Pending Sign Up
	admins, err := s.repo.FindAdmins()
	log.Println(admins)
	if err == nil {
		for _, admin := range admins {
			notification := &models.Notification{
				ID:      utils.GenerateUUID(),
				UserID:  admin.ID,
				Message: fmt.Sprintf("New sign up request from %s %s", newUser.FirstName, newUser.LastName),
				IsRead:  false,
			}

			_ = s.repo.CreateNotification(notification)
		}
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

func (s *Service) SendSignUpForm(req SendSignUpFormRequest) error {
	emailAddr := req.Email
	role := req.Role

	existingUser, err := s.repo.FindUserByEmail(emailAddr)
	if err != nil {
		return appErr.NewBadRequest("Failed to verify if the email exists", err)
	}

	if existingUser != nil {
		return appErr.NewBadRequest("User with that email already exists", nil)
	}

	token := utils.GenerateUUIDStr()

	invite := SendSignUpFormRequest{Email: emailAddr, Role: role}
	data, _ := json.Marshal(invite)

	err = s.rdb.Set(context.Background(), "signup_invite:"+token, data, 12*time.Hour).Err()
	if err != nil {
		return appErr.NewInternal("Failed to store sign-up invite", err)
	}

	signupURL := fmt.Sprintf("http://localhost:3000/owner&driver/signup?token=%s", token)

	if err := email.SendSignUpForm(emailAddr, role, signupURL); err != nil {
		return appErr.NewBadRequest("Invalid Email or User Role", err)
	}

	return nil
}

func (s *Service) SignUpDecision(req SignUpDecisionRequest, signUpID string) (*JWTAuthResponse, string, error) {
	ctx := context.Background()

	// Retrieving pending sign-up data from Redis
	val, err := s.rdb.Get(ctx, "pending_signup:"+signUpID).Result()
	if err == redis.Nil {
		return nil, "", appErr.NewBadRequest("Pending sign up not found", nil)
	} else if err != nil {
		return nil, "", appErr.NewInternal("Failed to retrieve pending sign up", err)
	}

	// Unmarshal the sign-up data
	var pendingSignUp PendingSignUp
	if err := json.Unmarshal([]byte(val), &pendingSignUp); err != nil {
		return nil, "", appErr.NewInternal("Failed to parse the pending sign up data", err)
	}

	// If rejected
	if !*req.IsAccepted {
		_ = s.rdb.Del(ctx, "pending_signup:"+signUpID).Err()
		return nil, "", nil
	}

	// If accepted
	createdUser, err := s.repo.CreateUser(pendingSignUp.User)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to create user at database", err)
	}

	_, err = s.repo.CreateAddress(pendingSignUp.Address)
	if err != nil {
		return nil, "", appErr.NewInternal("Failed to create address at database", err)
	}

	_ = s.rdb.Del(ctx, "pending_signup:"+signUpID).Err()

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
