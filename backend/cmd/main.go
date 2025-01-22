package main

import (
	"coachella-backend/config"              // Database configuration package
	_ "coachella-backend/docs"              // Swagger docs package (import for side effects)
	"coachella-backend/internal/handlers"   // Handlers
	"coachella-backend/internal/middleware" // Middleware for authentication and authorization
	"coachella-backend/internal/tasks"      // Scheduled tasks
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"       // Swagger UI files
	"github.com/swaggo/gin-swagger" // Gin Swagger middleware
	"log"
	"os"
	"path/filepath"
	"time"
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
	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize database connection
	config.ConnectDatabase()

	// Initialize scheduler for background tasks
	initializeScheduler()

	// Set up Gin router
	router := setupRouter()

	// Start the server
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// loadEnv loads the environment variables from the .env file
func loadEnv() error {
	wd, _ := os.Getwd()
	log.Println("Current working directory:", wd)

	err := godotenv.Load(filepath.Join("..", ".env"))
	if err != nil {
		return err
	}

	log.Println("SMTP Host:", os.Getenv("SMTP_HOST")) // Example log
	return nil
}

// initializeScheduler sets up and starts the task scheduler
func initializeScheduler() {
	scheduler := gocron.NewScheduler(time.Local)

	// Schedule tasks
	scheduler.Every(1).Day().At("09:00").Do(handlers.SendEventReminders) // Event reminders
	go func() {
		for {
			tasks.CleanUpExpiredTransactions() // Cleanup expired transactions
			time.Sleep(1 * time.Minute)
		}
	}()

	scheduler.StartAsync()
}

// setupRouter initializes the Gin router and routes
func setupRouter() *gin.Engine {
	r := gin.Default()

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Test endpoint for email
	r.GET("/test-email", handlers.SendTestEmail)

	// Authentication routes
	r.POST("/auth/admin-login", handlers.AdminLogin)
	r.POST("/auth/user-login", handlers.UserLogin)

	// Admin routes (protected)
	adminGroup := r.Group("/admin", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		adminGroup.POST("/tickets", handlers.CreateTicket)
		adminGroup.PUT("/tickets/:id", handlers.UpdateTicket)
		adminGroup.DELETE("/tickets/:id", handlers.DeleteTicket)
	}

	// User routes (protected)
	userGroup := r.Group("/user", middleware.AuthMiddleware(), middleware.RoleMiddleware("user"))
	{
		userGroup.GET("/transactions", handlers.GetUserTransactions)
		userGroup.POST("/transactions", handlers.CreateTransaction)
	}

	// Waitlist routes
	waitlistGroup := r.Group("/waitlist")
	{
		waitlistGroup.POST("", handlers.AddUserToWaitlist)
		waitlistGroup.GET("", handlers.GetWaitlist)
	}

	// Ticket routes (public)
	ticketGroup := r.Group("/tickets")
	{
		ticketGroup.GET("", handlers.GetTickets)
		ticketGroup.GET("/:id", handlers.GetTicketByID)
	}

	// Transaction routes (admin-only)
	transactionGroup := r.Group("/transactions", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		transactionGroup.GET("", handlers.GetTransactions)
		transactionGroup.GET("/:id", handlers.GetTransactionByID)
		transactionGroup.GET("/user-transactions", handlers.GetUserTransactions)
	}

	// Notification routes (user-specific)
	notificationGroup := r.Group("/notifications", middleware.AuthMiddleware())
	{
		notificationGroup.GET("", handlers.GetNotifications)
		notificationGroup.PATCH("/:id", handlers.MarkNotificationAsRead)
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

	return r
}
