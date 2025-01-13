package models

import "time"

type Ticket struct {
	TicketID          uint      `gorm:"primaryKey" json:"ticket_id"`
	EventID           uint      `gorm:"not null;index" json:"event_id" example:"1"`
	Event             Event     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Batch             int       `gorm:"not null" json:"batch" example:"1"`
	Type              string    `gorm:"type:varchar(255)" json:"type" example:"VIP"`
	Description       string    `gorm:"type:text" json:"description" example:"VIP access to the main stage."`
	Price             float64   `gorm:"type:decimal(10,2)" json:"price" example:"250.00"`
	QuantityAvailable int       `gorm:"not null" json:"quantity_available" example:"100"`
	// @swagger:ignore
	StartDate         DateOnly  `json:"start_date"`
	// @swagger:ignore
	EndDate           DateOnly  `json:"end_date"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	// Add replacement annotations for Swagger
	StartDateString string `json:"start_date" example:"15-01-2025" swaggerignore:"true"`
	EndDateString   string `json:"end_date" example:"28-02-2025" swaggerignore:"true"`
}


