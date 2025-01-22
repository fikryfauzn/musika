package tasks

import (
	"coachella-backend/config"
	"coachella-backend/internal/models"
	"log"
	"time"
)

func CleanUpExpiredTransactions() {
	var expiredTransactions []models.Transaction

	// Find expired transactions
	config.DB.Where("payment_status = ? AND timeout <= ?", "Pending", time.Now()).Find(&expiredTransactions)

	for _, transaction := range expiredTransactions {
		// Update transaction status to "Failed"
		transaction.PaymentStatus = "Failed"
		config.DB.Save(&transaction)

		// Restore ticket quantity
		var ticket models.Ticket
		if err := config.DB.First(&ticket, transaction.TicketID).Error; err == nil {
			ticket.QuantityAvailable += transaction.Quantity
			config.DB.Save(&ticket)
		}

		log.Printf("Processed expired transaction: %d\n", transaction.TransactionID)
	}
}

