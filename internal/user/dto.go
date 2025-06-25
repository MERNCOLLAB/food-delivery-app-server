package user

import "food-delivery-app-server/models"

// Update User & Update Profile Picture
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Bio   *string `json:"bio,omitempty"`
	Phone *string `json:"phone,omitempty"`
}

type UpdateUserResponse struct {
	Name  *string     `json:"name,omitempty"`
	Email *string     `json:"email,omitempty"`
	Bio   *string     `json:"bio,omitempty"`
	Phone *string     `json:"phone,omitempty"`
	Role  models.Role `json:"role,omitempty"`
}

type UpdateProfilePictureRequest struct {
	userId      string
	imageFile   interface{}
	imageHeader interface{}
}

func NewUpdateUserResponse(user *models.User) *UpdateUserResponse {
	return &UpdateUserResponse{
		Name:  &user.Name,
		Email: &user.Email,
		Bio:   &user.Bio,
		Phone: &user.Phone,
		Role:  user.Role,
	}
}

// Get All Users
type GetUserResponse struct {
	Name           *string     `json:"name,omitempty"`
	Email          *string     `json:"email,omitempty"`
	Bio            *string     `json:"bio,omitempty"`
	ProfilePicture *string     `json:"profilePicture,omitempty"`
	Phone          *string     `json:"phone,omitempty"`
	Role           models.Role `json:"role,omitempty"`
}

func NewGetUserResponse(user *models.User) GetUserResponse {
	return GetUserResponse{
		Name:           &user.Name,
		Email:          &user.Email,
		Bio:            &user.Bio,
		Phone:          &user.Phone,
		ProfilePicture: &user.ProfilePicture,
		Role:           user.Role,
	}
}
