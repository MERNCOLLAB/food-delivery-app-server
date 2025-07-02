package restaurant

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

func (r *Repository) CreateAddress(address *models.Address) (*models.Address, error) {
	if err := r.db.Create(address).Error; err != nil {
		return nil, err
	}
	return address, nil
}

func (r *Repository) GetAddressByRestaurantID(restaurantID uuid.UUID) (*models.Address, error) {
	var address models.Address
	if err := r.db.First(&address, "restaurant_id = ?", restaurantID).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *Repository) UpdateRestaurantAddressByRestaurantID(restaurantID uuid.UUID, newAddress string, lat, long float64) (*models.Address, error) {
	var address models.Address
	if err := r.db.First(&address, "restaurant_id = ?", restaurantID).Error; err != nil {
		return nil, err
	}

	address.Address = newAddress
	address.Latitude = lat
	address.Longitude = long

	if err := r.db.Save(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *Repository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetRestaurantByOwner(uid uuid.UUID) ([]models.Restaurant, error) {
	var restaurant []models.Restaurant
	if err := r.db.Find(&restaurant, "owner_id = ?", uid).Error; err != nil {
		return nil, err
	}
	return restaurant, nil
}

func (r *Repository) GetRestaurantByID(restoId uuid.UUID) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := r.db.First(&restaurant, "id = ?", restoId).Error; err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func (r *Repository) UpdateRestaurant(restaurant *models.Restaurant) error {
	return r.db.Save(restaurant).Error
}

func (r *Repository) DeleteRestaurant(restoId uuid.UUID) error {
	var restaurant models.Restaurant
	if err := r.db.Delete(restaurant, "id = ?", restoId).Error; err != nil {
		return err
	}
	return nil
}
