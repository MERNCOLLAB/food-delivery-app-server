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

func (r *Repository) FindResetCodeByUserId(userId string) (*models.PasswordReset, error) {
	var resetPw models.PasswordReset
	result := r.db.First(&resetPw, "user_id = ?", userId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &resetPw, nil
}

func (r *Repository) SaveResetCode(resetpw models.PasswordReset) error {
	if err := r.db.Save(&resetpw).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteResetCodeByID(id string) error {
	return r.db.Delete(&models.PasswordReset{}, "id = ?", id).Error
}

func (r *Repository) UpdateUserPassword(user *models.User) error {
	return r.db.Model(user).Update("password", user.Password).Error
}
