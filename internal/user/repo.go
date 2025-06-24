package user

import (
	"gorm.io/gorm"

	"food-delivery-app-server/models"

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

    if req.Name != nil {
        user.Name = *req.Name
    }
    if req.Email != nil {
        user.Email = *req.Email
    }
    if req.Bio != nil {
        user.Bio = *req.Bio
    }
    if req.Phone != nil {
        user.Phone = *req.Phone
    }

    if err := r.db.Save(&user).Error; err != nil {
        return nil, err
    }

    return &user, nil
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
