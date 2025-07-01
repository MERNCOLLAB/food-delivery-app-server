package restaurant

import (
	"food-delivery-app-server/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindRestaurantByName(name string) (*models.Restaurant, error) {
	var restaurant models.Restaurant

	if err := r.db.First(&restaurant, "name = ?", name).Error; err != nil {
		return nil, nil
	}
	return &restaurant, nil
}

func (r *Repository) CreateRestaurant(restaurantData *models.Restaurant) (*models.Restaurant, error) {
	if err := r.db.Create(restaurantData).Error; err != nil {
		return nil, err
	}
	return restaurantData, nil
}

func (r *Repository) GetRestaurantByOwner() {

}

func (r *Repository) UpdateRestaurant() {

}

func (r *Repository) DeleteRestaurant() {

}
