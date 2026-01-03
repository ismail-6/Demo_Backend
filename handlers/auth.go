package handlers

import (
	"learning-app-backend/database"
	"learning-app-backend/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Login - Simple userId based login (no password required)
func Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.LoginResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Trim whitespace and validate
	req.UserID = strings.TrimSpace(req.UserID)
	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, models.LoginResponse{
			Success: false,
			Message: "User ID cannot be empty",
		})
		return
	}

	// Check if user exists
	var user models.User
	result := database.DB.Where("user_id = ?", req.UserID).First(&user)

	if result.Error != nil {
		// User doesn't exist, create new user
		user = models.User{
			UserID:   req.UserID,
			Username: req.UserID, // Using userID as username for simplicity
		}

		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.LoginResponse{
				Success: false,
				Message: "Failed to create user",
			})
			return
		}
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Success: true,
		Message: "Login successful",
		User:    &user,
	})
}

// Logout - Simple logout endpoint
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

// GetUser - Get user details
func GetUser(c *gin.Context) {
	userID := c.Param("userId")

	var user models.User
	result := database.DB.Where("user_id = ?", userID).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}
