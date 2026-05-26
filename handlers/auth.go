package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"

	// "go-fiber-api/config"
	middleware "go-fiber-api/middlewares"
	"go-fiber-api/models"
	"go-fiber-api/utils"
)

// type AuthHandler struct {
// 	cfg *config.Config
// }

// func NewAuthHandler(cfg *config.Config) *AuthHandler {
// 	return &AuthHandler{cfg: cfg}
// }

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string          `json:"access_token"`
	User        models.AuthUser `json:"user"`
	Status      string          `json:"status"`
	Message     string          `json:"message"`
}

func (h *ApiHandler) LoginHandler(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to read request body: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid request"})
	}

	var user models.AuthUser
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid credentials"})
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		log.Printf("User %s wrong password", req.Username)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid credentials"})
	}

	if user.Status != "ACTIVE" {
		log.Printf("User %s is not active", req.Username)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid credentials"})
	}

	token, err := utils.CreateToken(user.UserId, user.OpId, user.Username, user.Email, user.Telephone, user.Level, user.RoleAction, h.cfg)
	if err != nil {
		log.Printf("User %s failed to generate token", req.Username)
		return c.JSON(fiber.Map{"status": "error", "message": "Failed to log in"})
	}

	// return c.JSON(fiber.Map{
	// 	"accessToken": token,
	// 	"user":        user,
	// 	"status":      "success",
	// 	"message":     "Welcome Back",
	// })

	return c.JSON(LoginResponse{
		AccessToken: token,
		User:        user,
		Status:      "success",
		Message:     "Welcome Back",
	})
}

func (h *ApiHandler) RefreshHandler(c *fiber.Ctx) error {
	accessToken := c.Locals("jwt").(string)
	if accessToken == "" {
		return c.JSON(fiber.Map{"status": "error", "message": "Failed to authentication"})
	}

	user := middleware.GetUserClaim(c)

	return c.JSON(LoginResponse{
		AccessToken: accessToken,
		Status:      "success",
		User:        user,
		Message:     "Welcome Back",
	})
}
