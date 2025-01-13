package email

import (
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
)

// SendEmail sends an email using SMTP configuration
func SendEmail(to, subject, body string) error {
	// Load SMTP details from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	emailFrom := os.Getenv("EMAIL_FROM")

	// Create a new mailer
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", emailFrom)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	// Convert SMTP_PORT to int
	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}

	// Set up the dialer
	dialer := gomail.NewDialer(smtpHost, port, smtpUsername, smtpPassword)

	// Send the email
	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("Failed to send email to %s: %v\n", to, err)
		return err
	}
	log.Printf("Email sent successfully to %s\n", to)
	return nil
}
