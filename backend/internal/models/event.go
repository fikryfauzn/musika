package models

import "time"


type Event struct {
	EventID         uint            `gorm:"primaryKey" json:"event_id"`
	Name            string          `gorm:"type:varchar(255)" json:"name"`
	Description     string          `gorm:"type:text" json:"description"`
	LocationCity    string          `gorm:"type:varchar(255)" json:"location_city"`
	LocationState   string          `gorm:"type:varchar(255)" json:"location_state"`
	LocationCountry string          `gorm:"type:varchar(255)" json:"location_country"`
	StartDate       DateOnly `json:"start_date"` // Indonesian date format
	EndDate         DateOnly `json:"end_date"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
