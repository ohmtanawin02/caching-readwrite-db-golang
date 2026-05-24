package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	domain "golang-fiber/internal/user/domain"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func ValidateLoginRequest(c *fiber.Ctx, v *validator.Validate) (*LoginRequest, error) {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return nil, err
	}
	if err := v.Struct(req); err != nil {
		return nil, err
	}
	return req, nil
}

func ToLoginDomainInput(req *LoginRequest) domain.LoginInput {
	return domain.LoginInput{
		Username: req.Username,
		Password: req.Password,
	}
}
