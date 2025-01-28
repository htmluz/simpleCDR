package controllers

import (
	"database/sql"
	"log"
	"radiusgo/models"
	"radiusgo/services"

	"github.com/gofiber/fiber/v2"
)

func HandleGetHomerCalls(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filters := new(models.HomerFilterParams)
		if err := c.QueryParser(filters); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":  "Parametros invalidos",
				"detail": err.Error(),
			})
		}
		calls, err := services.GetBriefHomerCalls(filters, db)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":  "Erro buscando chamadas no Homer",
				"detail": err.Error(),
			})
		}
		log.Println(calls)
		return c.JSON(fiber.Map{
			"calls": calls,
		})
	}
}

func HandleGetMessages(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		callID := c.Params("call_id")
		if callID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "call_id necessario pra busca",
			})
		}
		messages := services.GetHomerMessages(callID, db)
		return c.JSON(fiber.Map{
			"call_id":  callID,
			"messages": messages,
		})
	}
}
