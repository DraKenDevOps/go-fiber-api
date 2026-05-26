package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"go-fiber-api/config"
	"go-fiber-api/database"
	"go-fiber-api/handlers"
	middleware "go-fiber-api/middlewares"
	"go-fiber-api/routes"
	"go-fiber-api/zplogger"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := zplogger.InitLogger(cfg); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zplogger.SyncLogger()

	zplogger.Logger.Info("Application starting")

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(middleware.LogRequestResponse([]string{"/health", "/metrics"}))

	app.Get("/health", func(c *fiber.Ctx) error {
		res := fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"uptime":    time.Since(time.Unix(0, 0)) / time.Second,
			"version":   cfg.AppVersion,
		}
		return c.JSON(res)
	})

	app.Static("/static", filepath.Join(cfg.Cwd, "uploads"))

	// app.Use(cfg.BasePath + "/v1")
	routes.SetupRoutes(app, cfg, handlers.NewApiHandler(cfg, db))

	go func() {
		addr := cfg.Host + ":" + cfg.Port
		zplogger.Logger.Info("Server listening",
			zap.String("bound_on", addr),
		)

		if err := app.Listen(addr); err != nil {
			zplogger.Logger.Fatal("Server error",
				zap.Error(err),
			)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	zplogger.Logger.Info("Shutting down server")

	if err := app.ShutdownWithTimeout(30); err != nil {
		zplogger.Logger.Error("Server shutdown error",
			zap.Error(err),
		)
	}

	zplogger.Logger.Info("Server stopped")
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	requestID := c.Get("X-Request-ID")
	zplogger.Logger.Error(message,
		zap.String("requestId", requestID),
		zap.String("path", c.Path()),
		zap.String("method", c.Method()),
		zap.Int("statusCode", code),
		zap.Error(err),
	)

	return c.Status(code).JSON(fiber.Map{
		"status":    "error",
		"message":   message,
		"requestId": requestID,
	})
}
