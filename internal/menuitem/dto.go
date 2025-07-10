package menuitem

import (
	"food-delivery-app-server/models"
	"mime/multipart"
)

type CreateMenuItemRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	ImageFile   multipart.File
	ImageHeader *multipart.FileHeader
}

type CreateMenuItemResponse struct {
	ID           string  `json:"id"`
	RestaurantID string  `json:"restaurantId"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	Price        float64 `json:"price"`
	IsAvailable  bool    `json:"isAvailable"`
}

type UpdateMenuItemRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Category    *string  `json:"category,omitempty"`
	ImageFile   *multipart.File
	ImageHeader *multipart.FileHeader
}

type UpdateMenuItemResponse struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Category    *string  `json:"category,omitempty"`
	ImageURL    *string  `json:"imageURL,omitempty"`
}

type GetMenuItemByRestaurantResponse struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"imageURL"`
	IsAvailable bool    `json:"isAvailable"`
}

func NewGetMenuItemByRestoResponse(menuItems []models.MenuItem) []GetMenuItemByRestaurantResponse {
	var menuItemsList []GetMenuItemByRestaurantResponse
	for _, items := range menuItems {
		resp := GetMenuItemByRestaurantResponse{
			Name:        items.Name,
			Description: items.Description,
			Price:       items.Price,
			Category:    items.Category,
			ImageURL:    items.ImageURL,
			IsAvailable: items.IsAvailable,
		}
		menuItemsList = append(menuItemsList, resp)
	}
	return menuItemsList
}
