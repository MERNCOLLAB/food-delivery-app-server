package notifications

import (
	"food-delivery-app-server/models"
	"food-delivery-app-server/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type NotificationRepo interface {
	CreateNotification(notification *models.Notification) error
	GetRestaurantByID(id uuid.UUID) (*models.Restaurant, error)
}

func CreateStatusChangeNotifications(
	repo NotificationRepo,
	order *models.Order,
	newStatus models.Status,
) error {
	var notifications []*models.Notification

	switch newStatus {
	case models.AcceptedByOwner:
		// Notify customer that their order has been accepted
		if order.CustomerID != nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    *order.CustomerID,
				OrderID:   &order.ID,
				Message:   "Your order has been accepted and is being prepared!",
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}

	case models.RejectedByOwner:
		// Notify customer that their order has been rejected
		if order.CustomerID != nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    *order.CustomerID,
				OrderID:   &order.ID,
				Message:   "Your order has been rejected. ",
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}

	case models.ReadyForPickUp:
		// Notify customer that their order is ready for pickup
		if order.CustomerID != nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    *order.CustomerID,
				OrderID:   &order.ID,
				Message:   "Your order is ready for pickup! A driver will be assigned soon.",
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}

	case models.AcceptedByDriver:
		// Notify customer that a driver has accepted their order
		if order.CustomerID != nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    *order.CustomerID,
				OrderID:   &order.ID,
				Message:   "A driver has accepted your order and is on the way!",
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}

		// Notify restaurant owner that a driver has been assigned
		restaurant, err := repo.GetRestaurantByID(order.RestaurantID)
		if err == nil && restaurant.OwnerID != uuid.Nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    restaurant.OwnerID,
				OrderID:   &order.ID,
				Message:   "A driver has been assigned to pick up the order.",
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}

	case models.RejectedByDriver:
		// Notify customer that no driver accepted their order
		if order.CustomerID != nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    *order.CustomerID,
				OrderID:   &order.ID,
				Message:   "No driver was available to deliver your order. We'll try to find another driver.",
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}

		// Notify restaurant owner about driver rejection
		restaurant, err := repo.GetRestaurantByID(order.RestaurantID)
		if err == nil && restaurant.OwnerID != uuid.Nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    restaurant.OwnerID,
				OrderID:   &order.ID,
				Message:   "No driver accepted the order. It's back in the queue.",
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}

	case models.Assigned:
		// Notify customer that their order has been assigned to a driver
		if order.CustomerID != nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    *order.CustomerID,
				OrderID:   &order.ID,
				Message:   "Your order has been assigned to a driver and is being picked up!",
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}

	case models.InTransit:
		// Notify customer that their order is on the way
		if order.CustomerID != nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    *order.CustomerID,
				OrderID:   &order.ID,
				Message:   "Your order is on the way! The driver is heading to your location.",
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}

	case models.Delivered:
		// Notify customer that their order has been delivered
		if order.CustomerID != nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    *order.CustomerID,
				OrderID:   &order.ID,
				Message:   "Your order has been delivered! Enjoy your meal!",
				IsRead:    false,
				CreatedAt: time.Now(),
			})

		}

		// Notify restaurant owner that the order has been delivered
		restaurant, err := repo.GetRestaurantByID(order.RestaurantID)
		if err == nil && restaurant.OwnerID != uuid.Nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    restaurant.OwnerID,
				OrderID:   &order.ID,
				Message:   "Order has been successfully delivered to the customer.",
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}

	case models.Canceled:
		// Notify restaurant owner that the order has been canceled
		restaurant, err := repo.GetRestaurantByID(order.RestaurantID)
		if err == nil && restaurant.OwnerID != uuid.Nil {
			notifications = append(notifications, &models.Notification{
				ID:        utils.GenerateUUID(),
				UserID:    restaurant.OwnerID,
				OrderID:   &order.ID,
				Message:   "Order has been canceled by the customer.",
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}

	}

	for _, notification := range notifications {
		if err := repo.CreateNotification(notification); err != nil {
			return err
		}
	}

	return nil
}
