package restaurant

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

func (s *Service) CreateRestaurant(createRestaurantData CreateRestaurantRequest) (*CreateRestaurantResponse, error) {
	userId := createRestaurantData.UserId
	name := createRestaurantData.Name
	description := createRestaurantData.Description
	phone := createRestaurantData.Phone
	file := createRestaurantData.ImageFile
	fileHeader := createRestaurantData.ImageHeader

	if phone == "" || name == "" {
		return nil, appErr.NewBadRequest("Phone and Name is required", nil)
	}

	uid, err := utils.ParseId(userId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid ID", err)
	}

	foundRestaurantName, err := s.repo.FindRestaurantByName(name)
	if err != nil {
		return nil, appErr.NewInternal("Failed to check for existing restaurant name", err)
	}

	if foundRestaurantName != nil {
		return nil, appErr.NewBadRequest("Restaurant Name already exist", nil)
	}

	url, _, err := utils.UploadImage(file, fileHeader)
	if err != nil {
		return nil, appErr.NewInternal("Failed to upload the image", err)
	}
	restaurantID := utils.GenerateUUID()

	restaurantData := &models.Restaurant{
		ID:          restaurantID,
		OwnerID:     uid,
		Name:        name,
		Description: utils.SafeString(description, ""),
		Phone:       phone,
		ImageURL:    url,
	}

	newRestaurant, err := s.repo.CreateRestaurant(restaurantData)
	if err != nil {
		return nil, appErr.NewInternal("Failed to create the restaurant at the database", err)
	}

	filteredRestaurant := CreateRestaurantResponse{
		ID:      newRestaurant.ID.String(),
		OwnerID: newRestaurant.OwnerID.String(),
		Name:    newRestaurant.Name,
	}

	return &filteredRestaurant, nil
}

func (s *Service) GetRestaurantByOwner() {

}

func (s *Service) UpdateRestaurant() {

}

func (s *Service) DeleteRestaurant() {

}
