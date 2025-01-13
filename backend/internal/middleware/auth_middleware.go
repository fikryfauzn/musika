package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"os"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// AuthMiddleware checks the validity of a JWT token and sets claims in the context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Ensure the token starts with "Bearer "
		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header malformed"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix to extract the token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token's signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		// Handle token parsing errors
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Set claims in the context for downstream handlers
		if userID, ok := claims["id"]; ok {
			c.Set("id", userID)
		}
		if email, ok := claims["email"]; ok {
			c.Set("email", email)
		}
		if role, ok := claims["type"]; ok {
			c.Set("role", role)
		}

		// Proceed to the next handler
		c.Next()
	}
}

// RoleMiddleware restricts access based on the user's role
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the user's role from the context (set by AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role information missing"})
			c.Abort()
			return
		}

		// Check if the user's role matches the required role
		if role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}
