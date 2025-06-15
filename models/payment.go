package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID     `gorm:"primaryKey;type:uuid" json:"id"`
	CustomerID    uuid.UUID     `gorm:"type:uuid;index" json:"customerId"`
	Amount        float64       `gorm:"not null;" json:"amount"`
	PaymentMethod string        `gorm:"type:varchar(20); not null;" json:"paymentMethod"`
	PaidAt        time.Time     `gorm:"autoCreateTime" json:"paidAt"`
	PaymentStatus PaymentStatus `gorm:"type:varchar(20); not null" json:"status"`
	Customer      User          `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"customer"`
}
