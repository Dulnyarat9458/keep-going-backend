package models

import (
	"time"

	"gorm.io/gorm"
)

type HabitTracker struct {
	gorm.Model
	UserID        uint      `gorm:"not null"`
	Title         string    `gorm:"type:varchar(255);not null"`
	StartDate     time.Time `gorm:"not null"`
	LastResetDate time.Time `gorm:"not null"`
}
