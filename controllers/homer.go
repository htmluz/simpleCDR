package controllers

import (
	"radiusgo/models"
	"radiusgo/services"

	"github.com/gofiber/fiber/v2"
)

func HandleGetHomerCalls(c *fiber.Ctx) error {
	filters := new(models.FilterParamsHomer)
	if err := c.QueryParser(filters); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":  "Invalid parameters",
			"detail": err.Error(),
		})
	}
	results, err := services.QueryLoki(filters)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	calls, err := services.ProcessLokiResponse(results)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(calls)
}
