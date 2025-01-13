package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserID    uint           `gorm:"primaryKey" json:"user_id"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"uniqueIndex;size:255;not null" json:"email"` // Changed to VARCHAR(255)
	Password  string         `gorm:"not null" json:"-"`                          // Hashed password
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                             // Soft delete
}
