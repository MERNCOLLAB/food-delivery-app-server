package restaurant

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateRestaurant() {

}

func (r *Repository) GetRestaurantByOwner() {

}

func (r *Repository) UpdateRestaurant() {

}

func (r *Repository) DeleteRestaurant() {

}
