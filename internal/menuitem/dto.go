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
