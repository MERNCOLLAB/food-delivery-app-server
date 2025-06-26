package models

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID           uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
	UserID       *uuid.UUID `gorm:"type:uuid;index" json:"userId,omitempty"`
	RestaurantID *uuid.UUID `gorm:"type:uuid;index" json:"restaurantId,omitempty"`

	Address   string    `gorm:"type:varchar(100);not null" json:"address"`
	Label     string    `gorm:"type:varchar(10)" json:"label"`
	IsDefault bool      `json:"isDefault"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	User       *User       `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Restaurant *Restaurant `gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE" json:"restaurant,omitempty"`
}
