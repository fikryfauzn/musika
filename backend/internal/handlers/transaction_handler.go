package handlers

import (
	"coachella-backend/config"
	"coachella-backend/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetTransactions retrieves all transactions
// @Summary Retrieve all transactions
// @Description Get a list of all transactions, including user and ticket details
// @Tags Transactions
// @Produce json
// @Success 200 {array} models.Transaction
// @Failure 500 {object} models.GenericResponse "Internal Server Error"
// @Router /transactions [get]
func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction
	result := config.DB.
		Preload("User").
		Preload("Ticket").
		Preload("Ticket.Event"). // Preload the nested Event
		Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// GetTransactionByID retrieves a single transaction by ID
// @Summary Retrieve a transaction by ID
// @Description Get a transaction's details, including user and ticket details
// @Tags Transactions
// @Param id path int true "Transaction ID"
// @Produce json
// @Success 200 {object} models.Transaction
// @Failure 404 {object} models.GenericResponse "Transaction not found"
// @Router /transactions/{id} [get]
func GetTransactionByID(c *gin.Context) {
	id := c.Param("id")
	var transaction models.Transaction
	result := config.DB.
		Preload("User").
		Preload("Ticket").
		Preload("Ticket.Event"). // Preload the nested Event
		First(&transaction, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.GenericResponse{Error: "Transaction not found"})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

// GetUserTransactions retrieves transactions for a specific user
// @Summary Retrieve user transactions
// @Description Get all transactions for a specific user, including user and ticket details
// @Tags Transactions
// @Param user_id query int true "User ID"
// @Produce json
// @Success 200 {array} models.Transaction
// @Failure 400 {object} models.GenericResponse "Bad Request"
// @Failure 500 {object} models.GenericResponse "Internal Server Error"
// @Router /user-transactions [get]
func GetUserTransactions(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.GenericResponse{Error: "user_id is required"})
		return
	}

	var transactions []models.Transaction
	result := config.DB.
		Preload("User").         // Load User details
		Preload("Ticket").       // Load Ticket details
		Preload("Ticket.Event"). // Load nested Event details
		Where("user_id = ?", userID).
		Find(&transactions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
