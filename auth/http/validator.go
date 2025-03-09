package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Validator interface {
	Validate(c *fiber.Ctx) error
}


type BaseValidator struct {}

type AuthLocalValidator struct {
	BaseValidator
	Email string `json:"email" validate:"required,email"`
}

type AuthLocalVerificationValidator struct {
	BaseValidator
	Token string `json:"token" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func (b *BaseValidator) Validate(c *fiber.Ctx) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := c.BodyParser(b); err != nil {
		return err
	}

	if err := validate.Struct(b); err != nil {
		return err
	}

	return nil
}

func (a *AuthLocalValidator) Validate(c *fiber.Ctx) error {
	return a.BaseValidator.Validate(c)
}

func (a *AuthLocalVerificationValidator) Validate(c *fiber.Ctx) error {
	return a.BaseValidator.Validate(c)
}
