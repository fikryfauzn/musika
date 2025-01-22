package handlers

import (
    "coachella-backend/config"
    "coachella-backend/internal/models"
    "github.com/gin-gonic/gin"
    "net/http"
)

// AddUserToWaitlist adds a user to the waitlist for a sold-out ticket
// @Summary Add user to waitlist
// @Description Allows a user to join the waitlist for a sold-out ticket
// @Tags Waitlist
// @Accept json
// @Produce json
// @Param waitlist body models.Waitlist true "Waitlist Details"
// @Success 201 {object} models.Waitlist
// @Failure 400 {object} models.GenericResponse "Bad Request"
// @Failure 500 {object} models.GenericResponse "Internal Server Error"
// @Router /waitlist [post]
func AddUserToWaitlist(c *gin.Context) {
    var waitlist models.Waitlist
    if err := c.ShouldBindJSON(&waitlist); err != nil {
        c.JSON(http.StatusBadRequest, models.GenericResponse{Error: "Invalid input"})
        return
    }

    result := config.DB.Create(&waitlist)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: result.Error.Error()})
        return
    }

    c.JSON(http.StatusCreated, waitlist)
}

// GetWaitlist retrieves all users on the waitlist for a specific ticket
// @Summary Retrieve waitlist
// @Description Get a list of users on the waitlist for a specific ticket
// @Tags Waitlist
// @Param ticket_id query int true "Ticket ID"
// @Produce json
// @Success 200 {array} models.Waitlist
// @Failure 400 {object} models.GenericResponse "Bad Request"
// @Failure 500 {object} models.GenericResponse "Internal Server Error"
// @Router /waitlist [get]
func GetWaitlist(c *gin.Context) {
    ticketID := c.Query("ticket_id")
    if ticketID == "" {
        c.JSON(http.StatusBadRequest, models.GenericResponse{Error: "ticket_id is required"})
        return
    }

    var waitlist []models.Waitlist
    result := config.DB.Where("ticket_id = ?", ticketID).Preload("User").Preload("Ticket").Find(&waitlist)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, waitlist)
}
