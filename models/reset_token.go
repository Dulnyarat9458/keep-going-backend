package models

import (
	"time"

	"gorm.io/gorm"
)

type ResetToken struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	Token     string
	ExpiresAt time.Time
}
