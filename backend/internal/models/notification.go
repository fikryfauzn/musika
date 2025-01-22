package models

import "time"

type Notification struct {
	NotificationID  uint      `gorm:"primaryKey" json:"notification_id"`
	UserID          uint      `gorm:"not null" json:"user_id"` // Foreign key
	User            User      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	NotificationType string   `gorm:"type:varchar(50);not null" json:"notification_type"`
	Content         string    `gorm:"type:text" json:"content"`
	SentAt          *time.Time `json:"sent_at"` // Timestamp when sent
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
