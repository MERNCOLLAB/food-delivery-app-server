package resetpassword

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByEmail() {

}

func (r *Repository) SaveResetCode() {

}

func (r *Repository) UpdatePassword() {

}
