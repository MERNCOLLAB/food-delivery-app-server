package notification

import (
	"food-delivery-app-server/models"

	"github.com/google/uuid"

	appErr "food-delivery-app-server/pkg/errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserNotifications(userId string) ([]models.Notification, error) {
	uId, err := uuid.Parse(userId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid user ID", err)
	}

	notifications, err := s.repo.GetUserNotifications(uId)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (s *Service) MarkNotificationAsRead(notificationId, userId string) error {
	nId, err := uuid.Parse(notificationId)
	if err != nil {
		return appErr.NewBadRequest("Invalid notification ID", err)
	}

	uId, err := uuid.Parse(userId)
	if err != nil {
		return appErr.NewBadRequest("Invalid user ID", err)
	}

	err = s.repo.MarkNotificationAsRead(nId, uId)
	if err != nil {
		return err
	}
	return nil
}
