package models

import "time"

type Notification struct {
	NotificationID uint      `gorm:"primaryKey" json:"notification_id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`    // Foreign key
	User           User      `gorm:"constraint:OnDelete:CASCADE;"`    // Relationship
	Message        string    `gorm:"type:text" json:"message"`
	Type           string    `gorm:"type:enum('Reminder','Update','Confirmation')" json:"type"`
	IsRead         bool      `gorm:"default:false" json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}

