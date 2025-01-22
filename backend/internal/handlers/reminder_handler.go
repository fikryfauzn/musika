package handlers

import (
	"coachella-backend/config"
	"coachella-backend/internal/email"
	"coachella-backend/internal/models"
	"log"
	"path/filepath"
	"time"
)

// SendEventReminders checks for upcoming events and sends reminders to users
func SendEventReminders() {
	var transactions []models.Transaction

	// Find transactions for events happening tomorrow
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	result := config.DB.
		Preload("User").
		Preload("Ticket").
		Preload("Ticket.Event").
		Joins("JOIN tickets ON tickets.ticket_id = transactions.ticket_id").
		Joins("JOIN events ON events.event_id = tickets.event_id").
		Where("events.start_date = ?", tomorrow).
		Find(&transactions)

	if result.Error != nil {
		log.Printf("Error fetching transactions for reminders: %v\n", result.Error)
		return
	}

	for _, transaction := range transactions {
		// Prepare email content
		templateData := map[string]interface{}{
			"name":        transaction.User.Name,
			"event_name":  transaction.Ticket.Event.Name,
			"event_date":  transaction.Ticket.Event.StartDate.Format("02-01-2006"),
			"ticket_type": transaction.Ticket.Type,
			"quantity":    transaction.Quantity,
		}

		// Render email template
		templatePath := filepath.Join("..", "templates", "emails", "event_reminder.html")
		body, err := email.RenderTemplate(templatePath, templateData)
		if err != nil {
			log.Printf("Failed to render reminder email for user %d: %v\n", transaction.UserID, err)
			continue
		}

		// Send email
		err = email.SendEmail(transaction.User.Email, "Event Reminder: "+transaction.Ticket.Event.Name, body)
		if err != nil {
			log.Printf("Failed to send reminder email to user %d: %v\n", transaction.UserID, err)
		} else {
			log.Printf("Reminder email sent to user %d for event %s\n", transaction.UserID, transaction.Ticket.Event.Name)
		}
	}
}
