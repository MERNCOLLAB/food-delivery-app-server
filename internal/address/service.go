package address

import (
	"food-delivery-app-server/models"
	appErr "food-delivery-app-server/pkg/errors"
	"food-delivery-app-server/pkg/geocode"
	"food-delivery-app-server/pkg/utils"

	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAddress(req CreateAddressRequest, userId string) (*models.Address, error) {
	addr := req.Address

	uId, err := utils.ParseId(userId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid User ID", err)
	}

	existingAddr, err := s.repo.FindAddressByUser(addr, uId)
	if err != nil {
		return nil, err
	}
	if existingAddr != nil {
		return nil, appErr.NewBadRequest("Address already exists for this user", nil)
	}

	ctx := context.Background()
	lat, long, err := geocode.Geocode(ctx, addr)
	if err != nil {
		return nil, appErr.NewInternal("Failed to geocode the provided address", err)
	}

	addrId := utils.GenerateUUID()
	createAddr := &models.Address{
		ID:        addrId,
		UserID:    &uId,
		Address:   addr,
		Label:     req.Label,
		IsDefault: req.IsDefault,
		Latitude:  lat,
		Longitude: long,
	}

	newAddr, err := s.repo.CreateAddress(createAddr, uId)
	if err != nil {
		return nil, appErr.NewInternal("Failed to create address", err)
	}

	return newAddr, nil
}

func (s *Service) GetAddress(userId string) ([]models.Address, error) {
	uId, err := utils.ParseId(userId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid User ID", err)
	}

	addresses, err := s.repo.GetAddress(uId)
	if err != nil {
		return nil, appErr.NewInternal("Failed to get the addresses by user", err)
	}

	return addresses, nil
}

func (s *Service) UpdateAddress(addressId, userId string, req UpdateAddressRequest) (*models.Address, error) {
	addrId, err := utils.ParseId(addressId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid Address ID", err)
	}

	uId, err := utils.ParseId(userId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid User ID", err)
	}

	currentAddr, err := s.repo.FindAddressByIdAndUserId(addrId, uId)
	if err != nil {
		return nil, err
	}

	if currentAddr == nil {
		return nil, appErr.NewBadRequest("Address not found", nil)
	}

	if req.Address != nil && *req.Address != currentAddr.Address {
		existingAddr, err := s.repo.FindAddressByUser(*req.Address, uId)
		if err != nil {
			return nil, err
		}
		if existingAddr != nil {
			return nil, appErr.NewBadRequest("Address already exists for this user", nil)
		}

		ctx := context.Background()
		lat, long, err := geocode.Geocode(ctx, *req.Address)
		if err != nil {
			return nil, appErr.NewInternal("Failed to geocode the provided address", err)
		}

		currentAddr.Latitude = lat
		currentAddr.Longitude = long
	}

	if err := utils.Patch(currentAddr, &req); err != nil {
		return nil, appErr.NewInternal("Failed to patch address fields", err)
	}

	updatedAddr, err := s.repo.UpdateAddress(currentAddr)
	if err != nil {
		return nil, err
	}
	return updatedAddr, nil
}

func (s *Service) DeleteAddress(addressId, userId string) error {
	addrId, err := utils.ParseId(addressId)
	if err != nil {
		return appErr.NewBadRequest("Invalid Address ID", err)
	}

	uId, err := utils.ParseId(userId)
	if err != nil {
		return appErr.NewBadRequest("Invalid User ID", err)
	}

	currentAddr, err := s.repo.FindAddressByIdAndUserId(addrId, uId)
	if err != nil {
		return err
	}

	if currentAddr == nil {
		return appErr.NewBadRequest("Address not found", nil)
	}

	if currentAddr.IsDefault {
		return appErr.NewBadRequest("Cannot delete default address", nil)
	}

	allAddresses, err := s.repo.GetAddress(uId)
	if err != nil {
		return appErr.NewInternal("Failed to get user addresses", err)
	}

	if len(allAddresses) <= 1 {
		return appErr.NewBadRequest("Cannot delete the only address in user profile", nil)
	}

	err = s.repo.DeleteAddress(addrId, uId)
	if err != nil {
		return appErr.NewInternal("Failed to delete address", err)
	}

	return nil
}
