package utils

import (
	"time"

	"github.com/gistsapp/api/auth/config"
	"github.com/gofiber/fiber/v2"
)

func Cookie(key string, value string, config *config.CookiesConfig) *fiber.Cookie {
	cookie := new(fiber.Cookie)

	// will be overriden if auth is enabled
	cookie.Name = key
	// we are doing this because we don't want our user clients to have conflicting cookies
	if config.Auth.Enabled {
		if key == "access_token" {
			cookie.Name = config.Auth.AccessToken
		} else if key == "refresh_token" {
			cookie.Name = config.Auth.RefreshToken
		}
	}

	cookie.HTTPOnly = config.HTTPOnly
	cookie.Value = value
	cookie.Expires = time.Now().Add(time.Hour * 24 * 30 * 12) // 1 year
	cookie.Secure = config.Secure
	if config.Domain.Enabled {
		cookie.Domain = config.Domain.Value
	}
	return cookie
}

func ClearCookie(key string, env string, config *config.CookiesConfig) *fiber.Cookie {
	cookie := new(fiber.Cookie)

	cookie.Name = key
	if config.Auth.Enabled {
		if key == "access_token" {
			cookie.Name = config.Auth.AccessToken
		} else if key == "refresh_token" {
			cookie.Name = config.Auth.RefreshToken
		}
	}

	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour)
	cookie.Secure = config.Secure
	if config.Domain.Enabled {
		cookie.Domain = config.Domain.Value
	}

	return cookie
}
