package resetpassword

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

var user models.User

func (r *Repository) FindUserByEmail(email string) (*models.User, error) {
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) SaveResetCode(resetpw models.PasswordReset) error {
	if err := r.db.Save(&resetpw).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdatePassword() {

}
