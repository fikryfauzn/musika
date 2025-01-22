package models

import "time"

type Transaction struct {
	TransactionID uint      `gorm:"primaryKey" json:"transaction_id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`    // Foreign key
	User          User      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	TicketID      uint      `gorm:"not null;index" json:"ticket_id"` // Foreign key
	Ticket        Ticket    `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	TotalPrice    float64   `gorm:"type:decimal(10,2)" json:"total_price"`
	PaymentStatus string    `gorm:"type:enum('Pending','Paid','Failed','Expired')" json:"payment_status"`
	PaymentGateway string   `gorm:"type:varchar(255)" json:"payment_gateway"`
	Timeout       time.Time `gorm:"not null" json:"timeout"` // New field for timeout
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
