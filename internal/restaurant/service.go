package restaurant

import (
	"context"

	"food-delivery-app-server/models"
	appErr "food-delivery-app-server/pkg/errors"

	"food-delivery-app-server/pkg/geocode"
	"food-delivery-app-server/pkg/media"
	"food-delivery-app-server/pkg/sms"
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
	address := createReq.Address
	file := createReq.ImageFile
	fileHeader := createReq.ImageHeader

	// Data Validation and Preparations
	if phone == "" || name == "" || address == "" {
		return nil, appErr.NewBadRequest("Phone, Address, and Name is required", nil)
	}

	if err := sms.ValidatePhone(phone); err != nil {
		return nil, appErr.NewBadRequest("Invalid Phone Number Format", err)
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

	ctx := context.Background()
	lat, long, err := geocode.Geocode(ctx, address)
	if err != nil {
		return nil, appErr.NewInternal("Failed to geocode the provided address", err)
	}

	// Upload Image
	url, _, err := utils.UploadImage(file, fileHeader, "restaurants")
	if err != nil {
		return nil, appErr.NewInternal("Failed to upload the image", err)
	}

	// Restaurant Data Preparation
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

	// Address Data Preparation
	addressId := utils.GenerateUUID()
	newAddress := &models.Address{
		ID:           addressId,
		UserID:       &uid,
		RestaurantID: &restaurantID,
		Address:      address,
		IsDefault:    true,
		Latitude:     lat,
		Longitude:    long,
	}

	newAddr, err := s.repo.CreateAddress(newAddress)
	if err != nil {
		return nil, appErr.NewInternal("Failed to create address at the database", err)
	}

	filteredRestaurant := CreateRestaurantResponse{
		ID:      newRestaurant.ID.String(),
		OwnerID: newRestaurant.OwnerID.String(),
		Name:    newRestaurant.Name,
		Address: newAddr.Address,
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
		media.DeleteImage("restaurant", restaurant.ImageURL, "restaurants")

		url, _, err := utils.UploadImage(*updateReq.ImageFile, updateReq.ImageHeader, "restaurants")
		if err != nil {
			return nil, appErr.NewInternal("Failed to upload the image", err)
		}
		restaurant.ImageURL = url
	}

	if err := s.repo.UpdateRestaurant(restaurant); err != nil {
		return nil, appErr.NewInternal("Failed to update the restaurant", err)
	}

	var updatedAddress *string
	if updateReq.Address != nil && *updateReq.Address != "" {
		ctx := context.Background()
		lat, long, err := geocode.Geocode(ctx, *updateReq.Address)
		if err != nil {
			return nil, appErr.NewInternal("Failed to geocode the provided address", err)
		}

		addr, err := s.repo.UpdateRestaurantAddressByRestaurantID(restaurant.ID, *updateReq.Address, lat, long)
		if err != nil {
			return nil, appErr.NewInternal("Failed to update address at the database", err)
		}
		updatedAddress = &addr.Address
	} else {
		addr, _ := s.repo.GetAddressByRestaurantID(restaurant.ID)
		if addr != nil {
			updatedAddress = &addr.Address
		}
	}

	updatedResto := &UpdateRestaurantResponse{
		Name:        &restaurant.Name,
		Description: &restaurant.Description,
		Phone:       &restaurant.Phone,
		Address:     updatedAddress,
		ImageURL:    &restaurant.ImageURL,
	}

	return updatedResto, nil
}

func (s *Service) DeleteRestaurant(userId, restaurantId string) error {
	restoId, err := utils.ParseId(restaurantId)
	if err != nil {
		return appErr.NewBadRequest("Invalid Restaurant ID", err)
	}

	uid, err := utils.ParseId(userId)
	if err != nil {
		return appErr.NewBadRequest("Invalid User ID", err)
	}

	restaurant, err := s.repo.GetRestaurantByID(restoId)
	if err != nil {
		return appErr.NewNotFound("Restaurant not found", err)
	}

	if restaurant.OwnerID != uid {
		return appErr.NewUnauthorized("You are not authorized to delete this restaurant", nil)
	}

	media.DeleteImage("restaurant", restaurant.ImageURL, "restaurants")

	if err := s.repo.DeleteRestaurant(restoId); err != nil {
		return appErr.NewInternal("Failed to delete the restaurant", err)
	}

	return nil
}

func (s *Service) GetAllRestaurants() {

}

func (s *Service) GetMoreRestaurantDetails() {

}
