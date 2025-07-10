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

func (r *Repository) GetMenuItemByID(menuItemID uuid.UUID) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.First(&menuItem, "id = ?", menuItemID).Error; err != nil {
		return nil, err
	}
	return &menuItem, nil
}

func (r *Repository) GetRestaurantOwnerIDByID(restaurantID uuid.UUID) (uuid.UUID, error) {
	var ownerIDStr string
	err := r.db.
		Model(&models.Restaurant{}).
		Select("owner_id").
		Where("id = ?", restaurantID).
		Scan(&ownerIDStr).Error
	if err != nil {
		return uuid.Nil, err
	}
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		return uuid.Nil, err
	}
	return ownerID, nil
}

func (r *Repository) GetMenuItemByRestaurant(restoId uuid.UUID) ([]models.MenuItem, error) {
	var menuItems []models.MenuItem
	if err := r.db.Where("restaurant_id = ?", restoId).
		Find(&menuItems).
		Error; err != nil {
		return nil, err
	}
	return menuItems, nil
}

func (r *Repository) UpdateMenuItem(menuItem *models.MenuItem) error {
	return r.db.Save(menuItem).Error
}

func (r *Repository) DeleteMenuItem(menuId uuid.UUID) error {
	if err := r.db.Delete(&models.MenuItem{}, "id = ?", menuId).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetMenuItemDetails(menuID uuid.UUID) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.Where("id = ?", menuID).
		First(&menuItem).
		Error; err != nil {
		return nil, err
	}

	return &menuItem, nil
}
