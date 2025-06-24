package user

import "food-delivery-app-server/models"

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
    Email *string `json:"email,omitempty"`
    Bio   *string `json:"bio,omitempty"`
    Phone *string `json:"phone,omitempty"`
}

type UpdateUserResponse struct {
	Name  *string `json:"name,omitempty"`
    Email *string `json:"email,omitempty"`
    Bio   *string `json:"bio,omitempty"`
    Phone *string `json:"phone,omitempty"`
	Role models.Role `json:"role,omitempty"`
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
