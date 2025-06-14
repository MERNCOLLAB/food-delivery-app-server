package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name           string    `gorm:"type:varchar(100);not null" json:"name"`
	Email          string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password       string    `gorm:"not null" json:"password"`
	ProfilePicture string    `gorm:"type:text" json:"profilePicture"`
	Bio            string    `gorm:"type:text" json:"bio"`
	Phone          string    `gorm:"type:varchar(20)" json:"phone"`
	Role           Role      `gorm:"type:varchar(20)" json:"role"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
