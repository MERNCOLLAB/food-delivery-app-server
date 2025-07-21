package utils

import "food-delivery-app-server/models"

func IsValidOrderStatusTransition(currStatus, newStatus models.Status) bool {
	switch currStatus {
	case models.ReadyForPickUp:
		return newStatus == models.AcceptedByDriver || newStatus == models.RejectedByDriver
	case models.AcceptedByDriver:
		return newStatus == models.InTransit
	case models.InTransit:
		return newStatus == models.Delivered
	case models.Delivered:
		return false
	}

	return false
}
