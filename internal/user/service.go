package user

import (
	"mime/multipart"
	"regexp"

	appErr "food-delivery-app-server/pkg/errors"
	"food-delivery-app-server/pkg/media"
	"food-delivery-app-server/pkg/utils"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (s *Service) UpdateUser(req UpdateUserRequest, userId string) (*UpdateUserResponse, error) {
	uid, err := utils.ParseId(userId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid ID", err)
	}

	if req.Email != nil && !isValidEmail(*req.Email) {
		return nil, appErr.NewBadRequest("Invalid email format", nil)
	}

	updatedUser, err := s.repo.UpdateUser(uid, req)
	if err != nil {
		return nil, appErr.NewInternal("Failed to update the user", err)
	}

	return NewUpdateUserResponse(updatedUser), nil
}

func (s *Service) UpdateProfilePicture(updateProfilePicData UpdateProfilePictureRequest) (string, error) {
	userId := updateProfilePicData.userId
	file := updateProfilePicData.imageFile.(multipart.File)
	fileHeader := updateProfilePicData.imageHeader.(*multipart.FileHeader)

	uid, err := utils.ParseId(userId)
	if err != nil {
		return "", appErr.NewBadRequest("Invalid ID", err)
	}

	user, err := s.repo.FindUserByID(uid)
	if err != nil {
		return "", appErr.NewInternal("Failed to fetch the user", err)
	}

	media.DeleteProfilePicIfNotDefault(user.ProfilePicture, "profile_pictures")

	url, _, err := utils.UploadImage(file, fileHeader, "profile_pictures")
	if err != nil {
		return "", appErr.NewInternal("Failed to upload the image", err)
	}

	if err := s.repo.UpdateProfilePictureURL(uid, url); err != nil {
		return "", appErr.NewInternal("Failed to update the profile picture URL", err)
	}

	return url, nil
}

func (s *Service) DeleteUser(userId, email string) error {
	uid, err := utils.ParseId(userId)
	if err != nil {
		return appErr.NewBadRequest("Invalid ID", err)
	}

	user, err := s.repo.FindUserByID(uid)
	if err != nil {
		return appErr.NewNotFound("User not found", err)
	}

	if user.Email != email {
		return appErr.NewBadRequest("Email does not match", nil)
	}

	media.DeleteProfilePicIfNotDefault(user.ProfilePicture, "profile_pictures")

	if err := s.repo.DeleteUser(uid); err != nil {
		return appErr.NewInternal("Failed to delete the user", err)
	}

	return nil
}

func (s *Service) GetAllUsers() ([]GetUserResponse, error) {
	var userResponse []GetUserResponse

	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, appErr.NewInternal("Failed to query all users", err)
	}

	for _, user := range users {
		userResponse = append(userResponse, NewGetUserResponse(&user))
	}

	return userResponse, nil
}
