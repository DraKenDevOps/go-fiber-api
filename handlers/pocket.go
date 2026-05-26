package handlers

import (
	"go-fiber-api/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (h *ApiHandler) CreatePocket(c *fiber.Ctx) error {
	var body models.SavePocket

	if err := c.BodyParser(&body); err != nil {
		log.Printf("Failed to read request body: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	if len(body.Currency) != 3 {
		return c.JSON(fiber.Map{"status": "error", "message": "Currency must be 3 characters"})
	}

	if body.Currency == "" {
		body.Currency = "USD"
	}

	if body.Amount == nil {
		amnt := 0.00
		body.Amount = &amnt
	}

	var insertId int64
	err := h.db.Create(&body).Scan(&insertId).Error
	if err != nil {
		log.Printf("Failed to create pocket: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "Failed to create pocket"})
	}

	log.Printf("Success last pocket id: %d", insertId)
	return c.JSON(fiber.Map{"status": "success", "message": "User created successfully"})
}
