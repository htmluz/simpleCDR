package controllers

import (
	"database/sql"
	"radiusgo/models"

	"github.com/gofiber/fiber/v2"
)

func HandleUpdateCleanupDays(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.UpdateDaysRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":  "Invalid body",
				"detail": err.Error(),
			})
		}
		if req.Days <= 30 {
			return c.Status(400).JSON(fiber.Map{
				"error": "Value must be smaller than 30 days",
			})
		}
		q := `
		UPDATE cleanup_config
		SET cleanup_days = $1, updated_at = NOW()
		WHERE id = 1;
		`
		if _, err := db.Exec(q, req.Days); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":  "Error updating cleanup_days",
				"detail": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"message": "cleanup_days updated successfully",
			"days":    req.Days,
		})
	}
}

func HandleGetCleanupDays(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var res models.GetDaysResponse
		q := `SELECT cleanup_days, updated_at FROM cleanup_config LIMIT 1;`
		if err := db.QueryRow(q).Scan(&res.Days, &res.UpdatedAt); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":  "Error consulting the days",
				"detail": err.Error(),
			})
		}
		return c.JSON(res)
	}
}
