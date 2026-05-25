package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	domain "golang-fiber/internal/supplier/domain"
)

type CreateSupplierRequest struct {
	Name string `json:"name"        validate:"required"`
}

func ValidateCreateRequest(c *fiber.Ctx, v *validator.Validate) (*CreateSupplierRequest, error) {
	req := new(CreateSupplierRequest)
	if err := c.BodyParser(req); err != nil {
		return nil, err
	}
	if err := v.Struct(req); err != nil {
		return nil, err
	}
	return req, nil
}

func ToCreateDomainInput(req *CreateSupplierRequest) domain.CreateSupplierInput {
	return domain.CreateSupplierInput{
		Name: req.Name,
	}
}
