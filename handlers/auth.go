package handlers

import (
	"database/sql"
	"learning-app-backend/database"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type User struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Login - Query-driven login
func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}

	req.UserID = strings.TrimSpace(req.UserID)
	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User ID cannot be empty",
		})
		return
	}

	sqlDB, _ := database.DB.DB()

	// Check if user exists
	var user User
	query := `SELECT id, user_id, username, created_at, updated_at
			  FROM users WHERE user_id = $1 AND deleted_at IS NULL`

	err := sqlDB.QueryRow(query, req.UserID).Scan(
		&user.ID, &user.UserID, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create new user
		insertQuery := `INSERT INTO users (user_id, username, created_at, updated_at)
						VALUES ($1, $2, NOW(), NOW())
						RETURNING id, user_id, username, created_at, updated_at`

		err = sqlDB.QueryRow(insertQuery, req.UserID, req.UserID).Scan(
			&user.ID, &user.UserID, &user.Username, &user.CreatedAt, &user.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create user",
			})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"user":    user,
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

	sqlDB, _ := database.DB.DB()

	var user User
	query := `SELECT id, user_id, username, created_at, updated_at
			  FROM users WHERE user_id = $1 AND deleted_at IS NULL`

	err := sqlDB.QueryRow(query, userID).Scan(
		&user.ID, &user.UserID, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}
