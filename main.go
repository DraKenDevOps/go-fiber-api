package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"go-fiber-api/config"
	"go-fiber-api/database"
	"go-fiber-api/handlers"
	"go-fiber-api/routes"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	routes.SetupRoutes(app, cfg, handlers.NewApiHandler(cfg, db))

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(cfg.Host + ":" + cfg.Port))
}
