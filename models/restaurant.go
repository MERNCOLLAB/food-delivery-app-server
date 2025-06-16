package models

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	OwnerID     uuid.UUID `gorm:"type:uuid;index" json:"ownerId"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Phone       string    `gorm:"type:varchar(20)" json:"phone"`
	ImageURL    string    `gorm:"type:text" json:"imageURL"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Owner     *User      `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE" json:"owner,omitempty"`
	MenuItems []MenuItem `gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE" json:"menuItems,omitempty"`
}
