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
	ImageUrl     string    `gorm:"type:text" json:"imageURL"`
	IsAvailable  bool      `json:"isAvailable"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Restaurant Restaurant `gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE" json:"restaurant,omitempty"`
}
