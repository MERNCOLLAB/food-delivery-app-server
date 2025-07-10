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
