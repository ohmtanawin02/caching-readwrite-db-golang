package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	domain "golang-fiber/internal/product/domain"
	"golang-fiber/pkg/constants"
)

var ErrInvalidStatus = errors.New("status must be 'active' or 'inactive'")

type UpdateProductRequest struct {
	Name   string `json:"name"   validate:"required"`
	Price  float64 `json:"price"  validate:"required,gt=0"`
	Stock  int     `json:"stock"  validate:"gte=0"`
	Status string  `json:"status" validate:"required"`
}

func ValidateUpdateRequest(c *fiber.Ctx, v *validator.Validate) (*UpdateProductRequest, error) {
	req := new(UpdateProductRequest)
	if err := c.BodyParser(req); err != nil {
		return nil, err
	}
	if err := v.Struct(req); err != nil {
		return nil, err
	}
	if !constants.ProductStatus(req.Status).IsValid() {
		return nil, ErrInvalidStatus
	}
	return req, nil
}

func ToUpdateDomainInput(req *UpdateProductRequest) domain.UpdateProductInput {
	return domain.UpdateProductInput{
		Name:   req.Name,
		Price:  req.Price,
		Stock:  req.Stock,
		Status: constants.ProductStatus(req.Status),
	}
}
