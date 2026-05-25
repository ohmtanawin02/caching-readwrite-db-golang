package domain

import (
	"context"

	"golang-fiber/internal/supplier/domain/entity"
	"golang-fiber/pkg/constants"
)

type FindAllRequest struct {
	Page      int
	Limit     int
	SortBy    SupplierSortBy
	SortOrder constants.SortOrder
}

type FindAllResult struct {
	Items []entity.Supplier
	Total int64
}

type SupplierRepository interface {
	FindAll(context.Context, FindAllRequest) (FindAllResult, error)
	FindByID(context.Context, uint) (*entity.Supplier, error)
	FindByName(context.Context, string) (*entity.Supplier, error)
	Create(context.Context, *entity.Supplier) error
	Update(context.Context, *entity.Supplier) error
	SoftDelete(context.Context, uint) error
	Delete(context.Context, uint) error
}
