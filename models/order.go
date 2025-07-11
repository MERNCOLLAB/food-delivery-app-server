package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
	RestaurantID    uuid.UUID  `gorm:"type:uuid;index" json:"restaurantId"`
	CustomerID      *uuid.UUID `gorm:"type:uuid;index" json:"customerId"`
	DriverID        *uuid.UUID `gorm:"type:uuid;index" json:"driverId"`
	Status          Status     `gorm:"type:varchar(20);not null" json:"status"`
	TotalAmount     float64    `gorm:"not null" json:"totalAmount"`
	DeliveryFee     float64    `gorm:"not null" json:"deliveryFee"`
	DeliveryAddress string     `gorm:"type:varchar(100)" json:"deliveryAddress"`
	PlacedAt        time.Time  `gorm:"autoCreateTime" json:"placedAt"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`

	Payment    *Payment    `gorm:"foreignKey:OrderID" json:"payment,omitempty"`
	Restaurant *Restaurant `gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE" json:"restaurant,omitempty"`
	Customer   *User       `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	Driver     *User       `gorm:"foreignKey:DriverID;references:ID" json:"driver,omitempty"`

	OrderItems    []OrderItem    `gorm:"foreignKey:OrderID;" json:"items,omitempty"`
	Notifications []Notification `gorm:"foreignKey:OrderID" json:"notifications,omitempty"`
}
