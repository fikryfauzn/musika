package handlers

import (
	"coachella-backend/config"
	"coachella-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"os"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// AdminLogin allows admins to authenticate
// @Summary Admin Login
// @Description Authenticate an admin and issue a JWT token for accessing admin-specific routes
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Admin credentials (email and password)"
// @Success 200 {object} models.TokenResponse "JWT token for admin"
// @Failure 400 {object} models.GenericResponse "Invalid request body"
// @Failure 401 {object} models.GenericResponse "Invalid email or password"
// @Router /auth/admin-login [post]
func AdminLogin(c *gin.Context) {
	var credentials models.LoginRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, models.GenericResponse{Error: "Invalid request"})
		return
	}

	var admin models.Admin
	if err := config.DB.Where("email = ?", credentials.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.GenericResponse{Error: "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.GenericResponse{Error: "Invalid email or password"})
		return
	}

	token := generateJWT(admin.AdminID, admin.Email, "admin")
	c.JSON(http.StatusOK, models.TokenResponse{Token: token})
}

// UserLogin allows users to authenticate
// @Summary User Login
// @Description Authenticate a user and issue a JWT token for accessing user-specific routes
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User credentials (email and password)"
// @Success 200 {object} models.TokenResponse "JWT token for user"
// @Failure 400 {object} models.GenericResponse "Invalid request body"
// @Failure 401 {object} models.GenericResponse "Invalid email or password"
// @Router /auth/user-login [post]
func UserLogin(c *gin.Context) {
	var credentials models.LoginRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, models.GenericResponse{Error: "Invalid request"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ? AND deleted_at IS NULL", credentials.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.GenericResponse{Error: "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.GenericResponse{Error: "Invalid email or password"})
		return
	}

	token := generateJWT(user.UserID, user.Email, "user")
	c.JSON(http.StatusOK, models.TokenResponse{Token: token})
}

// Generate a JWT for authenticated users or admins
func generateJWT(id uint, email, userType string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"type":  userType, // 'user' or 'admin'
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}
