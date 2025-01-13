package models

import "time"


type Transaction struct {
	TransactionID uint      `gorm:"primaryKey" json:"transaction_id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`    // Foreign key
	User          User      `gorm:"constraint:OnDelete:CASCADE;"`    // Relationship
	TicketID      uint      `gorm:"not null;index" json:"ticket_id"` // Foreign key
	Ticket        Ticket    `gorm:"constraint:OnDelete:CASCADE;"`    // Relationship
	Quantity      int       `gorm:"not null" json:"quantity"`
	TotalPrice    float64   `gorm:"type:decimal(10,2)" json:"total_price"`
	PaymentStatus string    `gorm:"type:enum('Pending','Paid','Failed')" json:"payment_status"`
	PaymentGateway string   `gorm:"type:varchar(255)" json:"payment_gateway"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
