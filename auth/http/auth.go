package http

import (
	"github.com/gistsapp/api/auth/config"
	"github.com/gistsapp/api/auth/core"
	"github.com/gistsapp/api/auth/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Callback() fiber.Handler     //done
	Authenticate() fiber.Handler //done
	LocalAuth() fiber.Handler
	VerifyAuthToken() fiber.Handler
	Renew() fiber.Handler
	Logout() fiber.Handler
	Register(app *fiber.App)
}

type authController struct {
	service core.AuthService
	config  *config.Config
}


func NewAuthController(service core.AuthService, config *config.Config) AuthController {
	return authController{
		service: service,
		config:  config,
	}
}

// Callback godoc
//
//	@Summary		OAuth2 Callback
//	@Description	Use this endpoint to complete the OAuth2 flow
//	@Tags			auth
//	@Produce		json
//	@Sucess			302 {string} redirect to the client app
//	@Failure		404	{object}	http.HTTPErrorMessage
//	@Failure		400	{object}	http.HTTPErrorMessage
//	@Router			/auth/{provider}/callback [get]
func (a authController) Callback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := a.service.Callback(c)
		if err == core.UnkownProvider {
			return c.Status(fiber.ErrNotFound.Code).JSON(HTTPErrorMessage{
				Error: "Unkown provider",
			})
		}
		if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(HTTPErrorMessage{
				Error: "Couldn't complete auth",
			})
		}
		c.Cookie(utils.Cookie("access_token", token.AccessToken, &a.config.Cookies))
		c.Cookie(utils.Cookie("refresh_token", token.RefreshToken, &a.config.Cookies))
		return c.Redirect(a.config.Keycloak.RedirectURI)
	}
}

// Authenticate godoc
//
//	@Summary		Authenticate with redirect
//	@Description	Use this endpoint to authenticate with redirect
//	@Tags			auth
//	@Produce		json
//	@Sucess			302 {string} redirect to the client app
//	@Failure		400	{object}	http.HTTPErrorMessage
//	@Router			/auth/{provider} [get]
func (a authController) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return a.service.AuthenticateWithRedirect(c)
	}
}

// LocalAuth godoc
//
//	@Summary		Authenticate with code
//	@Description	Use this endpoint to authenticate with code
//	@Tags			auth
//	@Produce		json
//	@Sucess			200 {object} http.HTTPMessage
//	@Failure		400	{object}	http.HTTPErrorMessage
//	@Router			/auth/local/begin [post]
func (a authController) LocalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		e := new(AuthLocalValidator)
		if err := e.Validate(c); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(HTTPErrorMessage{
				Error: err.Error(),
			})
		}
		if _, err := a.service.AuthenticateWithCode(e.Email); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(HTTPErrorMessage{
				Error: err.Error(),
			})
		}

		return c.JSON(HTTPMessage{
			Message: "Hey check your emails :)",
		})
	}
}

// VerifyAuthToken godoc
//
//	@Summary		Verify auth token
//	@Description	Use this endpoint to verify auth token
//	@Tags			auth
//	@Produce		json
//	@Sucess			200 {object} http.HTTPTokens
//	@Failure		400	{object}	http.HTTPErrorMessage
//	@Router			/auth/local/verify [post]
func (a authController) VerifyAuthToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		e := new(AuthLocalVerificationValidator)

		if err := e.Validate(c); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(HTTPErrorMessage{
				Error: err.Error(),
			})
		}

		tokens, err := a.service.VerifyAuthToken(e.Token, e.Email)

		if err != nil {
			return c.Status(fiber.ErrUnauthorized.Code).JSON(HTTPErrorMessage{
				Error: err.Error(),
			})
		}

		c.Cookie(utils.Cookie("access_token", tokens.AccessToken, &a.config.Cookies))
		c.Cookie(utils.Cookie("refresh_token", tokens.RefreshToken, &a.config.Cookies))

		return c.JSON(
			HTTPTokens{
				AccessToken:  tokens.AccessToken,
				RefreshToken: tokens.RefreshToken,
			},
		)
	}
}

// Renew godoc
//
//	@Summary		Renew access token
//	@Description	Use this endpoint to renew access token
//	@Tags			auth
//	@Produce		json
//	@Sucess			200 {object} http.HTTPTokens
//	@Failure		400	{object}	http.HTTPErrorMessage
//	@Router			/auth/renew [get]
func (a authController) Renew() fiber.Handler {
	return func(c *fiber.Ctx) error {
		access_token := c.Locals("access_token").(string)

		tokens, err := a.service.Renew(access_token)

		if err != nil {
			return c.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Cookie(utils.Cookie("access_token", tokens.AccessToken, &a.config.Cookies))
		c.Cookie(utils.Cookie("refresh_token", tokens.RefreshToken, &a.config.Cookies))

		return c.JSON(HTTPTokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}

// Logout godoc
//
//	@Summary		Logout
//	@Description	Use this endpoint to logout (clear cookies)
//	@Tags			auth
//	@Produce		json
//	@Sucess			302 {string} redirect to the client app
//	@Failure		500	{object}	http.HTTPErrorMessage
//	@Router			/auth/logout [get]
func (a authController) Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Cookie(utils.ClearCookie("access_token", a.config.Keycloak.Realm, &a.config.Cookies))
		c.Cookie(utils.ClearCookie("refresh_token", a.config.Keycloak.Realm, &a.config.Cookies))
		return c.Redirect(a.config.Keycloak.RedirectURI)
	}
}

func (a authController) Register(app *fiber.App) {
	app.Get("/auth/:provider/callback", a.Callback())
	app.Get("/auth/:provider", a.Authenticate())
	app.Post("/auth/local/begin", a.LocalAuth())
	app.Post("/auth/local/verify", a.VerifyAuthToken())
	app.Get("/auth/renew", a.Renew())
	app.Get("/auth/logout", a.Logout())
}
