package auth

import (
	"context"
	"encoding/json"
	"fmt"
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

	if req.Token == "" {
		return nil, "", appErr.NewBadRequest("Missing invitation token", nil)
	}
	val, err := s.rdb.Get(context.Background(), "signup_invite:"+req.Token).Result()
	if err == redis.Nil {
		return nil, "", appErr.NewBadRequest("Invalid or expired invitation token", nil)
	} else if err != nil {
		return nil, "", appErr.NewInternal("Failed to verify invitation token", err)
	}

	var invite SendSignUpFormRequest
	if err := json.Unmarshal([]byte(val), &invite); err != nil {
		return nil, "", appErr.NewInternal("Failed to parse invitation data", err)
	}

	if invite.Email != req.Email {
		return nil, "", appErr.NewBadRequest("Invitation token does not match email", nil)
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

	stateID := utils.GenerateStateID(s.rdb, info)

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

	_, err := utils.GetTempUser(s.rdb, stateID)
	if err != nil {
		return appErr.NewBadRequest("Invalid or Expired State ID", err)
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
	stateID := req.StateID

	if err := sms.ValidatePhone(phone); err != nil {
		return nil, "", appErr.NewBadRequest("Invalid phone number", err)
	}

	storedOTP, err := utils.GetOTP(s.rdb, phone)
	if err != nil || storedOTP != otp {
		return nil, "", appErr.NewBadRequest("Invalid or expired OTP", nil)
	}

	oAuthData, err := utils.GetTempUser(s.rdb, stateID)
	if err != nil {
		return nil, "", appErr.NewBadRequest("Invalid or expired state ID", nil)
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
