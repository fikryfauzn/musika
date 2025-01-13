package handlers

import (
    "coachella-backend/config"
    "coachella-backend/internal/models"
    "github.com/gin-gonic/gin"
    "net/http"
)

// Get all users
func GetUsers(c *gin.Context) {
    var users []models.User
    result := config.DB.Find(&users)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

// Get user by ID
func GetUserByID(c *gin.Context) {
    id := c.Param("id")
    var user models.User
    result := config.DB.First(&user, id)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

// Create a new user
func CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result := config.DB.Create(&user)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
}

// Update an existing user
func UpdateUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User
    result := config.DB.First(&user, id)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    config.DB.Save(&user)
    c.JSON(http.StatusOK, user)
}

// Delete a user
func DeleteUser(c *gin.Context) {
    id := c.Param("id")
    result := config.DB.Delete(&models.User{}, id)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
