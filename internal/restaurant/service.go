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

func (s *Service) CreateRestaurant(userId string, createReq CreateRestaurantRequest) (*CreateRestaurantResponse, error) {
	name := createReq.Name
	description := createReq.Description
	phone := createReq.Phone
	file := createReq.ImageFile
	fileHeader := createReq.ImageHeader

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

func (s *Service) GetRestaurantByOwner(userId string) ([]GetRestaurantResponse, error) {
	uid, err := utils.ParseId(userId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid ID", err)
	}

	restaurantList, err := s.repo.GetRestaurantByOwner(uid)
	if err != nil {
		return nil, appErr.NewInternal("Failed to query restaurants by owner", err)
	}

	var formattedRestaurantList []GetRestaurantResponse
	for _, restaurant := range restaurantList {
		owner, _ := s.repo.GetUserByID(restaurant.OwnerID)
		resp := NewGetRestaurantResponse(&restaurant, owner)
		formattedRestaurantList = append(formattedRestaurantList, resp)
	}

	return formattedRestaurantList, nil
}

func (s *Service) UpdateRestaurant(restaurantId string, updateReq UpdateRestaurantRequest) (*UpdateRestaurantResponse, error) {
	restoId, err := utils.ParseId(restaurantId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid ID", err)
	}

	restaurant, err := s.repo.GetRestaurantByID(restoId)
	if err != nil {
		return nil, appErr.NewNotFound("Restaurant not found", err)
	}

	if err := utils.Patch(restaurant, &updateReq); err != nil {
		return nil, appErr.NewInternal("Failed to patch restaurant fields", err)
	}

	if updateReq.ImageFile != nil && updateReq.ImageHeader != nil {
		url, _, err := utils.UploadImage(*updateReq.ImageFile, updateReq.ImageHeader)
		if err != nil {
			return nil, appErr.NewInternal("Failed to upload the image", err)
		}
		restaurant.ImageURL = url
	}

	if err := s.repo.UpdateRestaurant(restaurant); err != nil {
		return nil, appErr.NewInternal("Failed to update the restaurant", err)
	}

	updatedResto := &UpdateRestaurantResponse{
		Name:        &restaurant.Name,
		Description: &restaurant.Description,
		Phone:       &restaurant.Phone,
		ImageURL:    &restaurant.ImageURL,
	}

	return updatedResto, nil
}

func (s *Service) DeleteRestaurant() {

}
