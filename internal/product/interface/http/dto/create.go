package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	domain "golang-fiber/internal/product/domain"
)

type CreateProductRequest struct {
	Name       string  `json:"name"        validate:"required"`
	Price      float64 `json:"price"       validate:"required,gt=0"`
	Stock      int     `json:"stock"       validate:"gte=0"`
	SupplierID *uint   `json:"supplier_id"`
}

func ValidateCreateRequest(c *fiber.Ctx, v *validator.Validate) (*CreateProductRequest, error) {
	req := new(CreateProductRequest)
	if err := c.BodyParser(req); err != nil {
		return nil, err
	}
	if err := v.Struct(req); err != nil {
		return nil, err
	}
	return req, nil
}

func ToCreateDomainInput(req *CreateProductRequest) domain.CreateProductInput {
	return domain.CreateProductInput{
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
		SupplierID: req.SupplierID,
	}
}
