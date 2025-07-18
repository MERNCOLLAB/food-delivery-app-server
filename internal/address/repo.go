package address

import (
	"food-delivery-app-server/models"
	appErr "food-delivery-app-server/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAddress(addr *models.Address, uId uuid.UUID) (*models.Address, error) {
	err := r.db.Create(addr).Error
	if err != nil {
		return nil, appErr.NewInternal("Failed to create address", err)
	}

	return addr, nil
}

func (r *Repository) FindAddressByUser(address string, uId uuid.UUID) (*models.Address, error) {
	var addr models.Address

	err := r.db.Where("user_id = ? AND address = ?", uId, address).First(&addr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, appErr.NewInternal("Failed to check existing address", err)
	}

	return &addr, nil
}

func (r *Repository) GetAddress(uId uuid.UUID) ([]models.Address, error) {
	var addresses []models.Address

	err := r.db.Where("user_id = ?", uId).Find(&addresses).Error
	if err != nil {
		return nil, appErr.NewInternal("Failed to get addresses", err)
	}
	return addresses, nil
}

func (r *Repository) UpdateAddress() {

}

func (r *Repository) DeleteAddress() {

}
