package menuitem

import (
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

func (s *Service) CreateMenuItem(restaurantId string, createReq CreateMenuItemRequest) (*CreateMenuItemResponse, error) {
	name := createReq.Name
	description := createReq.Description
	price := createReq.Price
	file := createReq.ImageFile
	fileHeader := createReq.ImageHeader

	if name == "" || price <= 0 {
		return nil, appErr.NewBadRequest("Missing required fields", nil)
	}

	restoId, err := utils.ParseId(restaurantId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid ID", err)
	}

	foundMenuItem, err := s.repo.FindMenuItemByName(name, restoId)
	if err != nil {
		return nil, appErr.NewInternal("Failed to check for existing menu item", err)
	}

	if foundMenuItem != nil {
		return nil, appErr.NewBadRequest("Menu Item already exist", nil)
	}

	url, _, err := utils.UploadImage(file, fileHeader, "menu-items")
	if err != nil {
		return nil, appErr.NewInternal("Failed to upload the image", err)
	}

	menuItemID := utils.GenerateUUID()

	menuItemData := &models.MenuItem{
		ID:           menuItemID,
		RestaurantID: restoId,
		Name:         name,
		Description:  utils.SafeString(description, ""),
		Price:        price,
		ImageURL:     url,
		IsAvailable:  true,
	}

	newMenuItem, err := s.repo.CreateMenuItem(menuItemData)
	if err != nil {
		return nil, appErr.NewInternal("Failed to create the menu item at the database", err)
	}

	filteredMenuItem := CreateMenuItemResponse{
		ID:           menuItemID.String(),
		RestaurantID: newMenuItem.RestaurantID.String(),
		Name:         newMenuItem.Name,
		Price:        newMenuItem.Price,
		IsAvailable:  newMenuItem.IsAvailable,
	}

	return &filteredMenuItem, nil
}

func (s *Service) GetMenuItemByRestaurant() {

}

func (s *Service) UpdateMenuItem() {
}

func (s *Service) DeleteMenuItem() {

}
