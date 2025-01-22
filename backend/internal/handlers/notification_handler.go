package handlers

import (
	"coachella-backend/config"
	"coachella-backend/internal/email"
	"coachella-backend/internal/models"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// GetNotifications retrieves all notifications for a user
// @Summary Retrieve notifications for a user
// @Description Get all notifications for a specific user
// @Tags Notifications
// @Param user_id query int true "User ID"
// @Produce json
// @Success 200 {array} models.Notification
// @Failure 400 {object} models.GenericResponse "Bad Request"
// @Failure 500 {object} models.GenericResponse "Internal Server Error"
// @Router /notifications [get]
func GetNotifications(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.GenericResponse{Error: "user_id is required"})
		return
	}

	var notifications []models.Notification
	result := config.DB.Where("user_id = ?", userID).Find(&notifications)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkNotificationAsRead marks a notification as read
// @Summary Mark a notification as read
// @Description Update a notification's status to "read"
// @Tags Notifications
// @Param id path int true "Notification ID"
// @Produce json
// @Success 200 {object} models.Notification
// @Failure 404 {object} models.GenericResponse "Notification not found"
// @Failure 500 {object} models.GenericResponse "Internal Server Error"
// @Router /notifications/{id} [patch]
func MarkNotificationAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	var notification models.Notification

	result := config.DB.First(&notification, notificationID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.GenericResponse{Error: "Notification not found"})
		return
	}

	config.DB.Save(&notification)
	c.JSON(http.StatusOK, notification)
}

func NotifyWaitlistedUsers(ticketID uint) {
	var waitlist []models.Waitlist
	result := config.DB.Where("ticket_id = ?", ticketID).Preload("User").Find(&waitlist)

	if result.Error != nil {
		log.Printf("Error fetching waitlist for ticket %d: %v\n", ticketID, result.Error)
		return
	}

	for _, entry := range waitlist {
		// Prepare notification content
		content := "Tickets are now available for the event you were waiting for!"

		notification := models.Notification{
			UserID:         entry.UserID,
			NotificationType: "Update", // Or "Reminder" if relevant
			Content:        content,
			CreatedAt:      time.Now(),
		}

		// Save notification to the database
		if err := config.DB.Create(&notification).Error; err != nil {
			log.Printf("Failed to create notification for user %d: %v\n", entry.UserID, err)
			continue
		}

		// Prepare email content
		templateData := map[string]interface{}{
			"name":       entry.User.Name,
			"ticket_id":  ticketID,
			"ticket_url": "http://example.com/tickets/" + string(ticketID), // Replace with actual URL
		}

		// Render email template
		templatePath := filepath.Join("..", "templates", "emails", "waitlist_notification.html")
		body, err := email.RenderTemplate(templatePath, templateData)
		if err != nil {
			log.Printf("Failed to render waitlist email for user %d: %v\n", entry.UserID, err)
			continue
		}

		// Send email
		err = email.SendEmail(entry.User.Email, "Tickets Now Available!", body)
		if err != nil {
			log.Printf("Failed to send waitlist email to user %d: %v\n", entry.UserID, err)
		} else {
			log.Printf("Waitlist email sent to user %d for ticket %d\n", entry.UserID, ticketID)
		}
	}

	// Optional: Clear waitlist entries for the notified ticket
	config.DB.Where("ticket_id = ?", ticketID).Delete(&models.Waitlist{})
}
