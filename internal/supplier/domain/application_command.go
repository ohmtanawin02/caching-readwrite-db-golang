package domain

import (
	"context"

	"golang-fiber/internal/supplier/domain/entity"
	"golang-fiber/pkg/constants"
)

type CreateSupplierInput struct {
	Name string
}

type UpdateSupplierInput struct {
	Name   string
	Status constants.SupplierStatus
}

type SupplierApplicationCommand interface {
	Create(context.Context, CreateSupplierInput) (*entity.Supplier, error)
	Update(context.Context, uint, UpdateSupplierInput) (*entity.Supplier, error)
	Delete(context.Context, uint) error
	SoftDelete(context.Context, uint) error
}
