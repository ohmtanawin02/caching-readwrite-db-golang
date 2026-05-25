package domain

import (
	"context"

	"golang-fiber/internal/product/domain/entity"
	"golang-fiber/pkg/constants"
)

type CreateProductInput struct {
	Name       string
	Price      float64
	Stock      int
	SupplierID *uint
}

type UpdateProductInput struct {
	Name   string
	Price  float64
	Stock  int
	Status constants.ProductStatus
}

type ProductApplicationCommand interface {
	Create(context.Context, CreateProductInput) (*entity.Product, error)
	Update(context.Context, uint, UpdateProductInput) (*entity.Product, error)
	Delete(context.Context, uint) error
	SoftDelete(context.Context, uint) error
}
