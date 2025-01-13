package handlers

import (
	"coachella-backend/config"
	"coachella-backend/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
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

	notification.IsRead = true
	config.DB.Save(&notification)
	c.JSON(http.StatusOK, notification)
}
