package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	Name           string    `gorm:"type:varchar(100);not null" json:"name"`
	Email          string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password       string    `gorm:"not null" json:"password"`
	ProfilePicture string    `gorm:"type:text" json:"profilePicture"`
	Bio            string    `gorm:"type:text" json:"bio"`
	Phone          string    `gorm:"type:varchar(20)" json:"phone"`
	Role           Role      `gorm:"type:varchar(20)" json:"role"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Addresses     []Address      `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	Restaurants   []Restaurant   `gorm:"foreignKey:OwnerID" json:"restaurants,omitempty"`
	Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
}
