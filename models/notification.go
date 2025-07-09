package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;index" json:"userId"`
	OrderID   *uuid.UUID `gorm:"type:uuid;index" json:"orderId,omitempty"`
	Message   string     `gorm:"type:varchar(50)" json:"message"`
	IsRead    bool       `json:"isRead"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`

	User  *User  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Order *Order `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
}
