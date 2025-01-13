package config

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "coachella-backend/internal/models"
)

var DB *gorm.DB

func ConnectDatabase() {
    dsn := "root:root@tcp(127.0.0.1:3306)/coachella?charset=utf8mb4&parseTime=True&loc=Local"
    database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database!")
    }

    // Migrate all models
	err = database.AutoMigrate(
        &models.User{},
        &models.Admin{},
        &models.Event{},
        &models.Ticket{},
        &models.Transaction{},
        &models.Notification{},
    )

    if err != nil {
        panic("Failed to migrate database!")
    }

    fmt.Println("Database connected and migrated successfully!")
    DB = database
}
