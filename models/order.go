package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	RestaurantID    uuid.UUID `gorm:"type:uuid;index" json:"restaurantId"`
	CustomerID      *uuid.UUID `gorm:"type:uuid;index" json:"customerId"`
	DriverID        *uuid.UUID `gorm:"type:uuid;index" json:"driverId"`
	Status          Status    `gorm:"type:varchar(20);not null" json:"status"`
	TotalAmount     float64   `gorm:"not null" json:"totalAmount"`
	DeliveryAddress string    `gorm:"type:varchar(100)" json:"deliveryAddress"`
	PlacedAt        time.Time `gorm:"autoCreateTime" json:"placedAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Restaurant Restaurant `gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE" json:"restaurant,omitempty"`
	Customer   User       `gorm:"foreignKey:CustomerID;constraint:OnDelete:SET NULL" json:"customer,omitempty"`
	Driver     *User       `gorm:"foreignKey:DriverID;constraint:OnDelete:SET NULL" json:"driver,omitempty"`

	OrderItems    []OrderItem    `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Notifications []Notification `gorm:"foreignKey:OrderID" json:"notifications,omitempty"`
}
