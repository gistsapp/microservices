package http

import (
    "github.com/gistsapp/api/auth/core"
    "github.com/gofiber/fiber/v2"
)

func JWTMiddleware(jwtService core.JWTService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        if token == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
        }

        claims, err := jwtService.VerifyAccessToken(token)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
        }

        c.Locals("userID", claims.UserID)
        return c.Next()
    }
}