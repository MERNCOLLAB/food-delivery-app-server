package address

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAddress() {

}

func (r *Repository) UpdateAddress() {

}

func (r *Repository) DeleteAddress() {

}
