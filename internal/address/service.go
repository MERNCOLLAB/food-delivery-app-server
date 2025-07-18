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

func (s *Service) GetAddress() {

}

func (s *Service) UpdateAddress() {

}

func (s *Service) DeleteAddress() {

}
