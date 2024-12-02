package controllers

import (
	"database/sql"
	"radiusgo/models"
	"radiusgo/utils"

	"github.com/gofiber/fiber/v2"
)

func HandleGetGateways(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var res []models.Gateway
		r, err := db.Query("SELECT name, ip FROM gateways")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":  "Database Error",
				"detail": err.Error(),
			})
		}
		defer r.Close()
		for r.Next() {
			var gw models.Gateway
			if err := r.Scan(&gw.Name, &gw.IP); err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": "Error parsing gateways data",
				})
			}
			res = append(res, gw)
		}
		return c.JSON(res)
	}
}

func HandlePostGateways(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.Gateway
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid Request",
			})
		}
		ip, er := utils.HasIP(req.IP)
		if er != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":  "Error parsing the IP address",
				"detail": er.Error(),
			})
		}

		_, err := db.Exec("INSERT INTO gateways (name, ip) VALUES ($1, $2)", req.Name, ip)
		if err != nil {
			return c.Status(409).JSON(fiber.Map{
				"error":  "Gateway already exists",
				"detail": err.Error(),
			})
		}
		return c.JSON(fiber.Map{"message": "Gateway created"})
	}
}

func HandleDeleteGateways(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.Gateway
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid Request",
			})
		}
		_, err := db.Exec("DELETE FROM gateways WHERE ip = $1", req.IP)
		if err != nil {
			return c.Status(409).JSON(fiber.Map{
				"error": "Gateway doesnt exists",
			})
		}
		return c.JSON(fiber.Map{"message": "Gateway deleted"})
	}
}
