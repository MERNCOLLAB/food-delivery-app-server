package order

import (
	"food-delivery-app-server/models"

	"github.com/google/uuid"
)

type PlaceOrderItem struct {
	MenuItem uuid.UUID `json:"menuItemId"`
	Quantity int       `json:"quantity"`
}

type PlaceOrderRequest struct {
	Items           []PlaceOrderItem `json:"items"`
	DeliveryAddress string           `json:"deliveryAddress"`
	AddressID       *uuid.UUID       `json:"addressId,omitempty"`
}

type PlaceOrderResponse struct {
	OrderID         uuid.UUID          `json:"orderId"`
	Status          models.Status      `json:"status"`
	TotalAmount     float64            `json:"totalAmount"`
	DeliveryFee     float64            `json:"deliveryFee"`
	DeliveryAddress string             `json:"deliveryAddress"`
	PlacedAt        string             `json:"placedAt"`
	Items           []models.OrderItem `json:"items"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

var allowedStatusTransitions = map[models.Status][]models.Status{
	models.Pending:          {models.AcceptedByOwner, models.RejectedByOwner, models.Canceled},
	models.AcceptedByOwner:  {models.ReadyForPickUp, models.RejectedByOwner, models.Canceled},
	models.ReadyForPickUp:   {models.AcceptedByDriver, models.RejectedByDriver},
	models.AcceptedByDriver: {models.Assigned},
	models.Assigned:         {models.InTransit},
	models.InTransit:        {models.Delivered},
	models.Delivered:        {},
	models.RejectedByOwner:  {},
	models.RejectedByDriver: {},
	models.Canceled:         {},
}
