package menuitem

import (
	"food-delivery-app-server/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindMenuItemByName(name string, restaurantId uuid.UUID) (*models.MenuItem, error) {
	var menuItem models.MenuItem

	if err := r.db.
		Where("name = ? AND restaurant_id = ?", name, restaurantId).
		First(&menuItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &menuItem, nil
}

func (r *Repository) CreateMenuItem(menuItemData *models.MenuItem) (*models.MenuItem, error) {
	if err := r.db.Create(menuItemData).Error; err != nil {
		return nil, err
	}
	return menuItemData, nil
}

func (r *Repository) GetMenuItemByRestaurant() {

}

func (r *Repository) UpdateMenuItem() {

}

func (r *Repository) DeleteMenuItem() {

}
