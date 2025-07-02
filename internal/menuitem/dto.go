package menuitem

import "mime/multipart"

type CreateMenuItemRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Price       float64 `json:"price"`
	ImageFile   multipart.File
	ImageHeader *multipart.FileHeader
}

type CreateMenuItemResponse struct {
	ID           string  `json:"id"`
	RestaurantID string  `json:"restaurantId"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	IsAvailable  bool    `json:"isAvailable"`
}

type UpdateMenuItemRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	ImageFile   *multipart.File
	ImageHeader *multipart.FileHeader
}

type UpdateMenuItemResponse struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	ImageURL    *string  `json:"imageURL,omitempty"`
}
