package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	AdminID   uint           `gorm:"primaryKey" json:"admin_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // Hashed password
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete, excluded from JSON
}
