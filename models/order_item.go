package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	OrderID    uuid.UUID `gorm:"type:uuid;index" json:"orderID"`
	MenuItemID uuid.UUID `gorm:"type:uuid;index" json:"menuItemId"`

	Price     float64   `gorm:"not null;" json:"price"`
	Quantity  int32     `gorm:"not null;" json:"quantity"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Order    *Order    `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
	MenuItem *MenuItem `gorm:"foreignKey:MenuItemID;references:ID" json:"menuItem,omitempty"`
}
