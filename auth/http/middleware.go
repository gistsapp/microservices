package http

import (
	"strings"

	"github.com/gistsapp/api/auth/core"
	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(jwtService core.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
		}

		bearer := strings.Split(header, " ")
		if len(bearer) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
		}
		token := bearer[1]
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
		}

		claims, err := jwtService.VerifyAccessToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("access_token", token)
		return c.Next()
	}
}

