package notification

import (
	"food-delivery-app-server/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserNotifications(userId uuid.UUID) ([]models.Notification, error) {
	var notifications []models.Notification

	err := r.db.Where("user_id = ?", userId).Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
