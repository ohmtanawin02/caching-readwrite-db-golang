package domain

import (
	"context"

	"golang-fiber/internal/supplier/domain/entity"
)

type SupplierApplicationQuery interface {
	FindAll(context.Context, FindAllRequest) (FindAllResult, error)
	FindByID(context.Context, uint) (*entity.Supplier, error)
}
