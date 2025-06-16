package models

import (
	"time"

	"github.com/google/uuid"
)

type PasswordReset struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ResetCode       string    `gorm:"type:varchar(255);not null;unique" json:"resetCode"`
	ResetCodeExpiry time.Time `gorm:"not null" json:"resetCodeExpiry"`
	IsUsed          bool      `gorm:"default:false" json:"isUsed"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
