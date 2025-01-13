package handlers

import (
    "coachella-backend/config"
    "coachella-backend/internal/models"
    "github.com/gin-gonic/gin"
    "net/http"
)

// Get event details
func GetEvent(c *gin.Context) {
    var event models.Event
    result := config.DB.First(&event)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
        return
    }
    c.JSON(http.StatusOK, event)
}

// Create or update the event (single event system)
func CreateOrUpdateEvent(c *gin.Context) {
    var event models.Event
    if err := c.ShouldBindJSON(&event); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if event exists
    var existingEvent models.Event
    result := config.DB.First(&existingEvent)
    if result.Error == nil {
        // Update the existing event
        event.EventID = existingEvent.EventID
        config.DB.Save(&event)
        c.JSON(http.StatusOK, event)
        return
    }

    // Create a new event
    config.DB.Create(&event)
    c.JSON(http.StatusCreated, event)
}

// Delete the event (optional, for admins)
func DeleteEvent(c *gin.Context) {
    result := config.DB.Delete(&models.Event{})
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
