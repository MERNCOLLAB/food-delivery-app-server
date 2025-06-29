package resetpassword

import (
	"food-delivery-app-server/models"
	"food-delivery-app-server/pkg/email"
	appErr "food-delivery-app-server/pkg/errors"
	"food-delivery-app-server/pkg/utils"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RequestResetPassword(req ResetPasswordRequest) error {
	emailAddr := req.Email

	user, err := s.repo.FindUserByEmail(emailAddr)
	if err != nil {
		return appErr.NewInternal("Failed to query email in database", err)
	}

	if user.Provider != "" {
		return appErr.NewBadRequest("Reset password feature is not available for OAuth account", nil)
	}

	code, expiresAt := email.GenerateResetCode()

	hashedCode, err := utils.HashPassword(code)
	if err != nil {
		return appErr.NewInternal("Failed to hash the generated code", err)
	}

	passwordReset := models.PasswordReset{
		ID:              utils.GenerateUUID(),
		UserID:          user.ID,
		ResetCode:       hashedCode,
		ResetCodeExpiry: expiresAt,
		IsUsed:          false,
	}

	err = s.repo.SaveResetCode(passwordReset)
	if err != nil {
		return appErr.NewInternal("Failed to save reset code", err)
	}

	err = email.SendResetCode(emailAddr, code)
	if err != nil {
		return appErr.NewInternal("Failed to send the reset code to the email", err)
	}

	return nil
}

func (s *Service) VerifyResetCode(req VerifyCodeRequest) error {
	emailAddr := req.Email
	resetCode := req.ResetCode

	user, err := s.repo.FindUserByEmail(emailAddr)
	if err != nil {
		return appErr.NewInternal("Failed to query email in database", err)
	}

	resetPw, err := s.repo.FindResetCodeByUserId(user.ID.String())
	if err != nil {
		return appErr.NewInternal("Failed to find the reset code", err)
	}
	if resetPw.IsUsed {
		return appErr.NewBadRequest("Reset code has already been used", nil)
	}
	if time.Now().After(resetPw.ResetCodeExpiry) {
		return appErr.NewBadRequest("Reset code has expired", nil)
	}
	if err := utils.ValidatePassword(resetPw.ResetCode, resetCode); err != nil {
		return appErr.NewBadRequest("Invalid reset code", nil)
	}

	if err := s.repo.DeleteResetCodeByID(resetPw.ID.String()); err != nil {
		return appErr.NewInternal("Failed to delete used reset code", err)
	}

	return nil
}

func (s *Service) UpdatePassword(req UpdatePasswordRequest) (*UpdatePasswordResponse, error) {
	email := req.Email
	password := req.Password
	confirmPassword := req.ConfirmPassword

	if password != confirmPassword {
		return nil, appErr.NewBadRequest("Password does not match", nil)
	}

	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, appErr.NewInternal("Failed to find user", err)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, appErr.NewInternal("Failed to hash password", err)
	}

	user.Password = hashedPassword
	if err := s.repo.UpdateUserPassword(user); err != nil {
		return nil, appErr.NewInternal("Failed to update password", err)
	}

	return &UpdatePasswordResponse{Email: user.Email}, nil
}
