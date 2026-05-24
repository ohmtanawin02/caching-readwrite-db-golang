package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	domain "golang-fiber/internal/user/domain"
)

type RegisterRequest struct {
	Username  string `json:"username"  validate:"required,min=3,max=100"`
	Email     string `json:"email"     validate:"required,email"`
	Password  string `json:"password"  validate:"required,min=8"`
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname"  validate:"required"`
	Phone     string `json:"phone"`
}

func ValidateRegisterRequest(c *fiber.Ctx, v *validator.Validate) (*RegisterRequest, error) {
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return nil, err
	}
	if err := v.Struct(req); err != nil {
		return nil, err
	}
	return req, nil
}

func ToRegisterDomainInput(req *RegisterRequest) domain.RegisterInput {
	return domain.RegisterInput{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Phone:     req.Phone,
	}
}
