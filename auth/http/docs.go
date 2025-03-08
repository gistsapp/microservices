package http

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v2"
)

type DocsHandler struct{}

func (d DocsHandler) Register(app *fiber.App) {
	app.Get("/docs", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			Theme: scalar.ThemeKepler,
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Gists auth service API",
			},
			DarkMode: true,
		})

		if err != nil {
			return c.Status(fiber.ErrInternalServerError.Code).JSON(HTTPErrorMessage{
				Error: err.Error(),
			})
		}

		return c.Format(htmlContent)
	})
}

func NewDocsHandler() DocsHandler {
	return DocsHandler{}
}
