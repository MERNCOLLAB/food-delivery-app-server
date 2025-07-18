package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (a Address) BeforeCreate(tx *gorm.DB) error {
	if a.IsDefault {
		return a.ensureSingleDefault(tx)
	}
	return nil
}

func (a Address) BeforeUpdate(tx *gorm.DB) error {
	if a.IsDefault {
		return a.ensureSingleDefault(tx)
	}
	return nil
}

func (a Address) ensureSingleDefault(tx *gorm.DB) error {
	var query *gorm.DB

	if a.UserID != nil {
		query = tx.Model(&Address{}).Where("user_id = ? AND id != ?", a.UserID, a.ID)
	} else if a.RestaurantID != nil {
		query = tx.Model(&Address{}).Where("restaurant_id = ? AND id != ?", a.RestaurantID, a.ID)
	} else {
		return nil
	}

	return query.Update("is_default", false).Error
}
