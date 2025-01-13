package handlers

import (
	"coachella-backend/config"
	"coachella-backend/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetTickets retrieves all tickets
// @Summary Retrieve all tickets
// @Description Get a list of all available tickets along with their associated event details
// @Tags Tickets
// @Produce json
// @Success 200 {array} models.Ticket "List of tickets"
// @Failure 500 {object} models.GenericResponse "Internal Server Error"
// @Router /tickets [get]
func GetTickets(c *gin.Context) {
	var tickets []models.Ticket
	// Use Preload to load the Event relationship
	result := config.DB.Preload("Event").Find(&tickets)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

// GetTicketByID retrieves a single ticket by ID
// @Summary Retrieve a ticket by ID
// @Description Get a single ticket's details using its ID, including event information
// @Tags Tickets
// @Param id path int true "Ticket ID"
// @Produce json
// @Success 200 {object} models.Ticket "Details of the ticket"
// @Failure 404 {object} models.GenericResponse "Ticket not found"
// @Router /tickets/{id} [get]
func GetTicketByID(c *gin.Context) {
	id := c.Param("id")
	var ticket models.Ticket

	// Use Preload to load the associated Event
	result := config.DB.Preload("Event").First(&ticket, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.GenericResponse{Error: "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// CreateTicket creates a new ticket
// @Summary Create a ticket
// @Description Add a new ticket to the system
// @Tags Tickets
// @Accept json
// @Produce json
// @Param ticket body models.Ticket true "Ticket details"
// @Success 201 {object} models.Ticket "The newly created ticket"
// @Failure 400 {object} models.GenericResponse "Bad Request"
// @Failure 500 {object} models.GenericResponse "Internal Server Error"
// @Router /tickets [post]
func CreateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, models.GenericResponse{Error: err.Error()})
		return
	}
	result := config.DB.Create(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, ticket)
}

// UpdateTicket updates an existing ticket
// @Summary Update a ticket
// @Description Modify the details of an existing ticket
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Param ticket body models.Ticket true "Updated ticket details"
// @Success 200 {object} models.Ticket "The updated ticket"
// @Failure 400 {object} models.GenericResponse "Bad Request"
// @Failure 404 {object} models.GenericResponse "Ticket not found"
// @Router /tickets/{id} [put]
func UpdateTicket(c *gin.Context) {
	id := c.Param("id")
	var ticket models.Ticket
	result := config.DB.First(&ticket, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.GenericResponse{Error: "Ticket not found"})
		return
	}
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, models.GenericResponse{Error: err.Error()})
		return
	}
	config.DB.Save(&ticket)
	c.JSON(http.StatusOK, ticket)
}

// DeleteTicket deletes a ticket
// @Summary Delete a ticket
// @Description Remove a ticket from the system
// @Tags Tickets
// @Param id path int true "Ticket ID"
// @Success 200 {object} models.GenericResponse "Ticket deleted successfully"
// @Failure 500 {object} models.GenericResponse "Internal Server Error"
// @Router /tickets/{id} [delete]
func DeleteTicket(c *gin.Context) {
	id := c.Param("id")
	result := config.DB.Delete(&models.Ticket{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.GenericResponse{
			Error: result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.GenericResponse{
		Message: "Ticket deleted successfully",
	})
}
