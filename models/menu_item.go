package models

import (
	"time"

	"github.com/google/uuid"
)

type MenuItem struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	RestaurantID uuid.UUID `gorm:"type:uuid;index" json:"restaurantId"`
	Name         string    `gorm:"type:varchar(100);not null;" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	Price        float64   `gorm:"not null;" json:"price"`
	Category     string    `gorm:"type:text" json:"category"`
	ImageURL     string    `gorm:"type:text" json:"imageURL"`
	IsAvailable  bool      `json:"isAvailable"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Restaurant *Restaurant `gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE" json:"restaurant,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:MenuItemID;constraint:OnDelete:CASCADE" json:"orderItems,omitempty"`
}
