package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID     `gorm:"primaryKey;type:uuid" json:"id"`
	OrderID       uuid.UUID     `gorm:"type:uuid;index" json:"orderId"`
	CustomerID    uuid.UUID     `gorm:"type:uuid;index" json:"customerId"`
	Amount        float64       `gorm:"not null;" json:"amount"`
	PaymentMethod string        `gorm:"type:varchar(20); not null;" json:"paymentMethod"`
	PaidAt        time.Time     `gorm:"autoCreateTime" json:"paidAt"`
	PaymentStatus PaymentStatus `gorm:"type:varchar(20); not null" json:"status"`

	Customer *User  `gorm:"foreignKey:CustomerID;references:ID;constraint:OnDelete:CASCADE" json:"customer,omitempty"`
	Order    *Order `gorm:"foreignKey:OrderID; references:ID; constraint:OnDelete:CASCADE" json:"order,omitempty"`
}
