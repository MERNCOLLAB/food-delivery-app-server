package user

import (
	"gorm.io/gorm"

	"food-delivery-app-server/models"
	"food-delivery-app-server/pkg/utils"

	"github.com/google/uuid"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

var user models.User

func (r *Repository) UpdateUser(uid uuid.UUID, req UpdateUserRequest) (*models.User, error) {
	if err := r.db.First(&user, "id = ?", uid).Error; err != nil {
		return nil, err
	}

	if err := utils.Patch(&user, &req); err != nil {
		return nil, err
	}

	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindUserByID(uid uuid.UUID) (*models.User, error) {
	if err := r.db.First(&user, "id = ?", uid).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateProfilePictureURL(uid uuid.UUID, url string) error {
	if err := r.db.First(&user, "id = ?", uid).Error; err != nil {
		return err
	}
	user.ProfilePicture = url

	return r.db.Save(&user).Error
}

func (r *Repository) DeleteUser(uid uuid.UUID) error {
	if err := r.db.Delete(&models.User{}, "id = ?", uid).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
