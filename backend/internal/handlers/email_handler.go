package handlers

import (
	"coachella-backend/internal/email"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SendTestEmail is a handler to test email functionality
func SendTestEmail(c *gin.Context) {
	// Replace with the recipient email, subject, and body content
	err := email.SendEmail(
		"ggzane23@gmail.com", // Replace with actual recipient
		"Test Email from Coachella",
		"<h1>Welcome!</h1><p>This is a test email from the Coachella system.</p>",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test email sent successfully!"})
}
