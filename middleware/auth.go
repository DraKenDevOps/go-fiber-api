package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"go-fiber-api/config"
	"go-fiber-api/models"
	"go-fiber-api/utils"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func AuthCheck(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {

		accessToken := c.Get("X-Access-Token")
		authorization := c.Get("Authorization")

		if authorization != "" {
			accessToken = strings.Replace(authorization, "Bearer ", "", 1)
		}
		if accessToken == "" {
			return c.JSON(Response{Status: "error", Message: "Authorization header required"})
		}

		claims, err := utils.Verify(accessToken, cfg)
		if err != nil {
			return c.JSON(Response{Status: "error", Message: "Invalid or expired token"})
		}

		c.Locals("jwt", accessToken)
		c.Locals("user_claims", claims)
		return c.Next()
	}
}

func GetUserClaim(c *fiber.Ctx) models.AuthUser {
	claims, ok := c.Locals("user_claims").(jwt.MapClaims)
	if !ok {
		return models.AuthUser{}
	}

	fmt.Printf("%+v\n", claims)
	fmt.Printf("%#v\n", claims)

	user := models.AuthUser{}

	if v, ok := claims["user_id"].(float64); ok {
		user.UserId = uint(v)
	}

	if v, ok := claims["username"].(string); ok {
		user.Username = v
	}

	if v, ok := claims["email"].(string); ok {
		user.Email = v
	}

	if v, ok := claims["telephone"].(string); ok {
		user.Telephone = v
	}

	if v, ok := claims["role_action"].(string); ok {
		user.RoleAction = models.UserRoleAction(v)
	}

	if v, ok := claims["level"].(string); ok {
		user.Level = models.UserLevel(v)
	}

	if v, ok := claims["op_id"].(float64); ok {
		opID := uint(v)
		user.OpId = &opID
	}

	return user
}
