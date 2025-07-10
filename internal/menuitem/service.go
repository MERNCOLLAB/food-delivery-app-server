package menuitem

import (
	"food-delivery-app-server/models"
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

func (s *Service) CreateMenuItem(restaurantId string, createReq CreateMenuItemRequest) (*CreateMenuItemResponse, error) {
	name := createReq.Name
	description := createReq.Description
	price := createReq.Price
	category := createReq.Category
	file := createReq.ImageFile
	fileHeader := createReq.ImageHeader

	if name == "" || category == "" || price <= 0 {
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
		Category:     category,
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
		Category:     newMenuItem.Category,
		IsAvailable:  newMenuItem.IsAvailable,
	}

	return &filteredMenuItem, nil
}

func (s *Service) GetMenuItemByRestaurant(restaurantId string) ([]GetMenuItemByRestaurantResponse, error) {
	restoId, err := utils.ParseId(restaurantId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid ID", err)
	}

	menuItems, err := s.repo.GetMenuItemByRestaurant(restoId)
	if err != nil {
		return nil, appErr.NewInternal("Failed to query menu items by restaurant", err)
	}

	return NewGetMenuItemByRestoResponse(menuItems), nil
}

func (s *Service) UpdateMenuItem(menuItemId string, updateReq UpdateMenuItemRequest) (*UpdateMenuItemResponse, error) {
	menuId, err := utils.ParseId(menuItemId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid ID", err)
	}

	menuItem, err := s.repo.GetMenuItemByID(menuId)
	if err != nil {
		return nil, appErr.NewNotFound("Menu item not found", err)
	}

	if err := utils.Patch(menuItem, &updateReq); err != nil {
		return nil, appErr.NewInternal("Failed to patch menu item fields", err)
	}

	if updateReq.ImageFile != nil && updateReq.ImageHeader != nil {
		media.DeleteImage("menu item", menuItem.ImageURL, "menu-items")
		url, _, err := utils.UploadImage(*updateReq.ImageFile, updateReq.ImageHeader, "menu-items")
		if err != nil {
			return nil, appErr.NewInternal("Failed to upload the image", err)
		}
		menuItem.ImageURL = url
	}

	if err := s.repo.UpdateMenuItem(menuItem); err != nil {
		return nil, appErr.NewInternal("Failed to update the menu item", err)
	}

	updatedMenu := &UpdateMenuItemResponse{
		Name:        &menuItem.Name,
		Description: &menuItem.Description,
		Price:       &menuItem.Price,
		Category:    &menuItem.Category,
		ImageURL:    &menuItem.ImageURL,
	}

	return updatedMenu, nil
}

func (s *Service) DeleteMenuItem(userId, menuItemId string) error {
	menuId, err := utils.ParseId(menuItemId)
	if err != nil {
		return appErr.NewBadRequest("Invalid Menu Item ID", err)
	}

	uid, err := utils.ParseId(userId)
	if err != nil {
		return appErr.NewBadRequest("Invalid User ID", err)
	}

	menuItem, err := s.repo.GetMenuItemByID(menuId)
	if err != nil {
		return appErr.NewNotFound("Menu Item not found", err)
	}

	restoId := menuItem.RestaurantID

	ownerID, err := s.repo.GetRestaurantOwnerIDByID(restoId)
	if err != nil {
		return appErr.NewNotFound("Menu Item not found", err)
	}
	if ownerID != uid {
		return appErr.NewUnauthorized("You are not authorized to delete this restaurant", nil)
	}

	media.DeleteImage("menu item", menuItem.ImageURL, "menu-items")

	if err := s.repo.DeleteMenuItem(menuId); err != nil {
		return appErr.NewInternal("Failed to delete the menu item", err)
	}

	return nil
}

func (s *Service) GetMoreMenuItemDetails(menuItemId string) (*models.MenuItem, error) {
	menuID, err := utils.ParseId(menuItemId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid Menu Item ID", err)
	}

	menuItem, err := s.repo.GetMenuItemDetails(menuID)
	if err != nil {
		return nil, appErr.NewInternal("Failed to query the menu item details", err)
	}

	return menuItem, nil
}
