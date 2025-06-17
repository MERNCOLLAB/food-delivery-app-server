package user

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) UpdateUser() {

}

func (r *Repository) FindUserByEmail() {

}

func (r *Repository) UpdateProfilePictureURL() {

}

func (r *Repository) FindUserByName() {

}

func (r *Repository) DeleteUser() {

}

func (r *Repository) GetAllUsers() {
}
