package controllers

import (
	"database/sql"
	"radiusgo/models"
	"radiusgo/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req models.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":  "Invalid request",
				"detail": err.Error(),
			})
		}

		var userID int
		var passwordHash string
		var role string
		err := db.QueryRow("SELECT id, password_hash, role FROM users WHERE username = $1", req.Username).
			Scan(&userID, &passwordHash, &role)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		accessToken, err := services.GenerateAccessToken(req.Username, role)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Error generating access_token",
			})
		}
		refreshToken, err := services.GenerateRefreshToken(req.Username)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Error generating refresh_token",
			})
		}

		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		_, err = db.Exec("INSERT into refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", userID, refreshToken, expiresAt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Error inserting token into the database",
			})
		}

		return c.JSON(fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func HandleRefreshToken(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.RefreshRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}
		var userID int
		var expiresAt time.Time
		err := db.QueryRow("SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1", req.RefreshToken).
			Scan(&userID, &expiresAt)
		if err == sql.ErrNoRows {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid refresh_token",
			})
		} else if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Database error",
			})
		}
		if time.Now().After(expiresAt) {
			return c.Status(401).JSON(fiber.Map{
				"error": "Expired refresh_token",
			})
		}
		expiresAt = time.Now().Add(-1 * time.Minute)
		db.QueryRow("UPDATE refresh_tokens SET expires_at = $1 WHERE token = $2", expiresAt, req.RefreshToken)

		var userRole string
		var username string
		err = db.QueryRow("SELECT username, role FROM users WHERE id = $1", userID).Scan(&username, &userRole)
		if err == sql.ErrNoRows {
			return c.Status(500).JSON(fiber.Map{
				"error": "User not found",
			})
		} else if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Database error",
				"details": err.Error(),
			})
		}
		newAccessToken, err := services.GenerateAccessToken(username, userRole)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Error generating new access_token",
			})
		}
		newRefreshToken, err := services.GenerateRefreshToken(username)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Error generating new refresh_token",
			})
		}
		expiresAt = time.Now().Add(7 * 24 * time.Hour)
		_, err = db.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", userID, newRefreshToken, expiresAt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Database error",
				"details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
		})
	}
}
