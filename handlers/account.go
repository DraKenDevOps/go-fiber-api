package handlers

import (
	"go-fiber-api/models"
	"log"

	"github.com/gofiber/fiber/v2"
	// "go-fiber-api/config"
	// "gorm.io/gorm"
)

// type AccountHandler struct {
// 	db  *gorm.DB
// 	cfg *config.Config
// }

// func NewAccountHandler(db *gorm.DB, cfg *config.Config) *AccountHandler {
// 	return &AccountHandler{db: db, cfg: cfg}
// }

func (h *ApiHandler) CreateAccount(c *fiber.Ctx) error {
	var account models.AccReqBody

	if err := c.BodyParser(&account); err != nil {
		log.Printf("Failed to read request body: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	var lastInsertId int64
	err := h.db.Create(&account).Scan(&lastInsertId).Error
	if err != nil {
		log.Printf("Failed to create account: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "Failed to create account"})
	}

	log.Printf("Success last account id: %d", lastInsertId)
	return c.JSON(fiber.Map{"status": "success", "message": "User created successfully"})
}
