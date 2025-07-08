package auth

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

func (r *Repository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) CreateAddress(address *models.Address) (*models.Address, error) {
	if err := r.db.Create(address).Error; err != nil {
		return nil, err
	}
	return address, nil
}

func (r *Repository) FindFacebookUserByProfilePicturePrefix(string) (*models.User, error) {
	prefix := "https://platform-lookaside.fbsbx.com"
	var user models.User
	if err := r.db.Where("profile_picture LIKE ? AND provider = ?", prefix+"%", "facebook").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindAdmins() ([]models.User, error) {
	var admins []models.User
	err := r.db.Where("role = ?", models.Admin).Find(&admins).Error
	return admins, err
}

func (r *Repository) CreateNotification(notification *models.Notification) error {
	return r.db.Create(notification).Error
}
