package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	domain "golang-fiber/internal/supplier/domain"
	"golang-fiber/pkg/constants"
)

var ErrInvalidStatus = errors.New("status must be 'active' or 'inactive'")

type UpdateSupplierRequest struct {
	Name   string `json:"name"   validate:"required"`
	Status string `json:"status" validate:"required"`
}

func ValidateUpdateRequest(c *fiber.Ctx, v *validator.Validate) (*UpdateSupplierRequest, error) {
	req := new(UpdateSupplierRequest)
	if err := c.BodyParser(req); err != nil {
		return nil, err
	}
	if err := v.Struct(req); err != nil {
		return nil, err
	}
	if !constants.SupplierStatus(req.Status).IsValid() {
		return nil, ErrInvalidStatus
	}
	return req, nil
}

func ToUpdateDomainInput(req *UpdateSupplierRequest) domain.UpdateSupplierInput {
	return domain.UpdateSupplierInput{
		Name:   req.Name,
		Status: constants.SupplierStatus(req.Status),
	}
}
