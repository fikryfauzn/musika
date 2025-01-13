package main

import (
	_ "coachella-backend/docs" // Swagger docs package (import for side effects)
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // Swagger UI files
	"github.com/swaggo/gin-swagger" // Gin Swagger middleware

	"coachella-backend/config"              // Database configuration package
	"coachella-backend/internal/handlers"   // Handlers
	"coachella-backend/internal/middleware" // Middleware for authentication and authorization
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"

)

// @title Coachella API Documentation
// @version 1.0
// @description This is the API documentation for the Coachella backend API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@coachella.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	wd, _ := os.Getwd()
	log.Println("Current working directory:", wd)


	// Explicitly load .env from the backend folder
	err := godotenv.Load(filepath.Join("..", ".env"))
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Example: Accessing the SMTP host
	smtpHost := os.Getenv("SMTP_HOST")
	log.Println("SMTP Host:", smtpHost)

	// Initialize database connection
	config.ConnectDatabase()

	// Set up Gin router
	r := gin.Default()

	r.GET("/test-email", handlers.SendTestEmail)

	// Register Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Authentication routes
	r.POST("/auth/admin-login", handlers.AdminLogin) // Admin login
	r.POST("/auth/user-login", handlers.UserLogin)   // User login

	// Admin routes (protected by AuthMiddleware and RoleMiddleware)
	adminGroup := r.Group("/admin", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		adminGroup.POST("/tickets", handlers.CreateTicket)       // Admin can create tickets
		adminGroup.PUT("/tickets/:id", handlers.UpdateTicket)    // Admin can update tickets
		adminGroup.DELETE("/tickets/:id", handlers.DeleteTicket) // Admin can delete tickets
	}

	// User routes (protected by AuthMiddleware and RoleMiddleware)
	userGroup := r.Group("/user", middleware.AuthMiddleware(), middleware.RoleMiddleware("user"))
	{
		userGroup.GET("/transactions", handlers.GetUserTransactions) // Users can view their transactions
	}

	// Public ticket routes (accessible to all)
	ticketGroup := r.Group("/tickets")
	{
		ticketGroup.GET("", handlers.GetTickets)        // List all tickets
		ticketGroup.GET("/:id", handlers.GetTicketByID) // Get ticket by ID
	}

	// Transaction routes (admin-only for general management)
	transactionGroup := r.Group("/transactions", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		transactionGroup.GET("", handlers.GetTransactions)                       // Admin can view all transactions
		transactionGroup.GET("/:id", handlers.GetTransactionByID)                // Admin can view specific transaction
		transactionGroup.GET("/user-transactions", handlers.GetUserTransactions) // Users can view their own transactions
	}

	// Notification routes (user-specific)
	notificationGroup := r.Group("/notifications", middleware.AuthMiddleware())
	{
		notificationGroup.GET("", handlers.GetNotifications)             // Users can view notifications
		notificationGroup.PATCH("/:id", handlers.MarkNotificationAsRead) // Mark notification as read
	}

	// Protected test route
	r.GET("/protected", middleware.AuthMiddleware(), func(c *gin.Context) {
		userID := c.GetString("id")
		userType := c.GetString("type")
		email := c.GetString("email")
		c.JSON(200, gin.H{
			"message":   "Access granted",
			"user_id":   userID,
			"user_type": userType,
			"email":     email,
		})
	})

	// Start the server
	r.Run(":8080")
}
