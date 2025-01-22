package models

import "time"

// Waitlist represents a record of a user waiting for a specific ticket type
type Waitlist struct {
    WaitlistID uint      `gorm:"primaryKey" json:"waitlist_id"`
    UserID     uint      `gorm:"not null;index" json:"user_id"`    // Foreign key
    User       User      `gorm:"constraint:OnDelete:CASCADE;"`    // Relationship to User
    TicketID   uint      `gorm:"not null;index" json:"ticket_id"` // Foreign key
    Ticket     Ticket    `gorm:"constraint:OnDelete:CASCADE;"`    // Relationship to Ticket
    CreatedAt  time.Time `json:"created_at"`
}
