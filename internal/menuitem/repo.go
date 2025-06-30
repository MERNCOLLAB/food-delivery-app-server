package menuitem

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateMenuItem() {

}

func (r *Repository) GetMenuItemByRestaurant() {

}

func (r *Repository) UpdateMenuItem() {

}

func (r *Repository) DeleteMenuItem() {

}
