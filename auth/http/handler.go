package http

import (
	"github.com/gistsapp/api/auth/core"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Register(app *fiber.App)
}


type HTTPTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type HTTPErrorMessage struct {
	Error string `json:"error"`
}

type HTTPMessage struct {
	Message string `json:"message"`
}
type handler struct {
	jwtService core.JWTService
}

func (h *handler) Register(app *fiber.App) {
	panic("unimplemented")
}

func NewHandler(jwtService core.JWTService) Handler {
	return &handler{
		jwtService: jwtService,
	}
}
