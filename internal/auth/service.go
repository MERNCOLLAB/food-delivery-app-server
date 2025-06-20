package auth

import (
	"regexp"

	"food-delivery-app-server/models"
	appErr "food-delivery-app-server/pkg/errors"
	"food-delivery-app-server/pkg/utils"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SignUp(req SignUpRequest) (*SignUpResponse, string, error) {
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
		return nil, "", appErr.NewBadRequest("Email already exists", nil)
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
		ProfilePicture: "https://fallback-image.com",
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

	userResponse := SignUpResponse{
		ID:   createdUser.ID.String(),
		Role: string(createdUser.Role),
	}

	return &userResponse, token, nil
}

func (s *Service) SignIn() {
}

func (s *Service) OAuth() {
}

func (s *Service) SignOut() {
}
