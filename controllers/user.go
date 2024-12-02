package controllers

import (
	"database/sql"
	"radiusgo/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func HandleGetUsers(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var res []models.User
		r, err := db.Query("SELECT username, role FROM users")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":  "Database error",
				"datils": err.Error(),
			})
		}
		defer r.Close()
		for r.Next() {
			var user models.User
			if err := r.Scan(&user.User, &user.Role); err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": "Error parsing user data",
				})
			}
			res = append(res, user)
		}
		return c.JSON(res)
	}
}

func HandleRegister(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":  "Invalid Request",
				"detail": err.Error(),
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Error generating password hash",
			})
		}
		_, err = db.Exec("INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3)",
			req.Username, string(hashedPassword), req.Role)
		if err != nil {
			return c.Status(409).JSON(fiber.Map{
				"error": "User already exists",
			})
		}

		return c.JSON(fiber.Map{"message": "User created"})
	}
}

func HandlePasswordChange(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		type PasswordChangeRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var req PasswordChangeRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid Request",
			})
		}
		userRole := c.Locals("role").(string)
		userName := c.Locals("username").(string)
		if userRole != "admin" && userName != req.Username {
			return c.Status(403).JSON(fiber.Map{
				"error": "You can only change your own password",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Error generating password",
			})
		}
		_, err = db.Exec("UPDATE users SET password_hash = $1 WHERE username = $2", hashedPassword, req.Username)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Error updating password",
			})
		}
		return c.JSON(fiber.Map{"message": "Password updated successfully"})
	}
}
