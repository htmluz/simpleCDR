package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"radiusgo/models"
	"radiusgo/services"

	"github.com/gofiber/fiber/v2"
)

func InsertIntoQ(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := c.Body()
		bilhete := new(models.Bilhete)

		if err := json.Unmarshal(data, bilhete); err != nil {
			log.Fatal("Erro fazendo parsing do json")
		}
		services.QueueHandleBilhete(db, bilhete)
		return c.SendStatus(201)
	}
}

func GetQ(c *fiber.Ctx) error {
	bilhetes, count := services.GetActiveCalls()
	return c.Status(200).JSON(fiber.Map{
		"message": bilhetes,
		"count":   count,
	})
}
