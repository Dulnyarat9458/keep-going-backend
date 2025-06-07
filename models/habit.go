package models

import (
	"time"

	"gorm.io/gorm"
)

type HabitTracker struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	Title         string         `json:"title" gorm:"type:varchar(255);not null"`
	StartDate     time.Time      `json:"start_date" gorm:"not null"`
	LastResetDate time.Time      `json:"last_reset_date" gorm:"not null"`
}
