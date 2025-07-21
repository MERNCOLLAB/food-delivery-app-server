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

func (r *Repository) MarkNotificationAsRead(nId, uId uuid.UUID) error {
	res := r.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", nId, uId).
		Update("is_read", true)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
