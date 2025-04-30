package models

import (
	"time"

	"gorm.io/gorm"
)

type HabitTracker struct {
	gorm.Model
	UserID        uint      `gorm:"not null"`
	Title         string    `json:"title" gorm:"type:varchar(255);not null"`
	StartDate     time.Time `json:"start_date" gorm:"not null"`
	LastResetDate time.Time `json:"last_reset_date" gorm:"not null"`
}
