package handlers

import (
	"coachella-backend/config"
	"coachella-backend/internal/email"
	"coachella-backend/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func ptr(t time.Time) *time.Time {
	return &t
}


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

// CreateTransaction creates a new transaction and sends a confirmation email
// @Summary Create a transaction
// @Description Create a new transaction for ticket purchase, send a confirmation email, and generate a notification
// @Tags Transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction Details"
// @Success 201 {object} models.Transaction
// @Failure 400 {object} models.GenericResponse "Bad Request"
// @Failure 404 {object} models.GenericResponse "Not Found"
// @Failure 500 {object} models.GenericResponse "Internal Server Error"
// @Router /user/transactions [post]
func CreateTransaction(c *gin.Context) {
	var transaction models.Transaction

	// Bind JSON payload
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, models.GenericResponse{Error: "Invalid input"})
		return
	}

	// Set default transaction values
	transaction.PaymentStatus = "Pending"
	transaction.PaymentGateway = "Midtrans" // Placeholder for future integration
	transaction.Timeout = time.Now().Add(15 * time.Minute)

	// Fetch the associated ticket
	var ticket models.Ticket
	if err := config.DB.Preload("Event").First(&ticket, transaction.TicketID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.GenericResponse{Error: "Ticket not found"})
		return
	}

	// Check if enough tickets are available
	if ticket.QuantityAvailable < transaction.Quantity {
		c.JSON(http.StatusBadRequest, models.GenericResponse{Error: "Not enough tickets available"})
		return
	}

	// Decrement ticket quantity
	ticket.QuantityAvailable -= transaction.Quantity
	if err := config.DB.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: "Failed to update ticket quantity"})
		return
	}

	// Save the transaction in the database
	if err := config.DB.Create(&transaction).Error; err != nil {
		// Rollback ticket quantity if transaction fails
		ticket.QuantityAvailable += transaction.Quantity
		config.DB.Save(&ticket)
		c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: "Failed to create transaction"})
		return
	}

	// Fetch the user details for email and notification
	var user models.User
	if err := config.DB.First(&user, transaction.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.GenericResponse{Error: "User not found"})
		return
	}

	// Prepare email content
	templateData := map[string]interface{}{
		"name":        user.Name,
		"event_name":  ticket.Event.Name,
		"ticket_type": ticket.Type,
		"quantity":    transaction.Quantity,
		"total_price": strconv.FormatFloat(transaction.TotalPrice, 'f', 2, 64),
		"event_date":  ticket.Event.StartDate.Format("02-01-2006"),
	}

	// Render the email template
	templatePath := filepath.Join("..", "templates", "emails", "purchase_confirmation.html")
	body, err := email.RenderTemplate(templatePath, templateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: "Failed to render email template"})
		return
	}

	// Send confirmation email
	err = email.SendEmail(user.Email, "Ticket Purchase Confirmation", body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: "Failed to send email"})
		return
	}

	// Save a notification to the database
	notification := models.Notification{
		UserID:          user.UserID,
		NotificationType: "Confirmation",
		Content:         "Your purchase for " + ticket.Event.Name + " (" + ticket.Type + ") has been confirmed.",
		SentAt:          ptr(time.Now()),
	}
	if err := config.DB.Create(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.GenericResponse{Error: "Failed to create notification"})
		return
	}

	// Respond with the created transaction
	c.JSON(http.StatusCreated, transaction)
}
