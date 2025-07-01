package restaurant

import (
	"food-delivery-app-server/models"
	"mime/multipart"
)

type CreateRestaurantRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Phone       string  `json:"phone"`
	ImageFile   multipart.File
	ImageHeader *multipart.FileHeader
}

type CreateRestaurantResponse struct {
	ID      string `json:"restaurantID"`
	OwnerID string `json:"userID"`
	Name    string `json:"name"`
}

type UpdateRestaurantRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	ImageFile   *multipart.File
	ImageHeader *multipart.FileHeader
}

type UpdateRestaurantResponse struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	ImageURL    *string `json:"imageURL,omitempty"`
}

type GetRestaurantResponse struct {
	Name           *string `json:"name"`
	OwnerFirstName *string `json:"ownerFirstName"`
	OwnerLastName  *string `json:"ownerLastName"`
	Description    *string `json:"description,omitempty"`
	Phone          *string `json:"phone"`
	ImageURL       *string `json:"imageURL"`
}

func NewGetRestaurantResponse(restaurant *models.Restaurant, owner *models.User) GetRestaurantResponse {
	var firstName, lastName *string
	if owner != nil {
		firstName = &owner.FirstName
		lastName = &owner.LastName
	}
	return GetRestaurantResponse{
		Name:           &restaurant.Name,
		OwnerFirstName: firstName,
		OwnerLastName:  lastName,
		Description:    &restaurant.Description,
		Phone:          &restaurant.Phone,
		ImageURL:       &restaurant.ImageURL,
	}
}
