package http

import "github.com/gofiber/fiber/v2"

type Handler interface {
	Register(app *fiber.App)
}
