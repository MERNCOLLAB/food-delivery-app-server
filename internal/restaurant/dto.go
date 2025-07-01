package restaurant

import "mime/multipart"

type CreateRestaurantRequest struct {
	UserId      string
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
