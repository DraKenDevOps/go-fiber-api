package routes

import (
	"github.com/gofiber/fiber/v2"

	"go-fiber-api/config"
	"go-fiber-api/handlers"
	"go-fiber-api/middleware"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, handler *handlers.ApiHandler) {
	api := app.Group("/api")
	api.Post("/login", handler.LoginHandler)
	api.Use(middleware.AuthCheck(cfg)).Get("/refresh", handler.RefreshHandler)
}
