package domain

import (
	"context"

	"golang-fiber/internal/product/domain/entity"
	"golang-fiber/pkg/constants"
)

type FindAllRequest struct {
	Page       int
	Limit      int
	SupplierID uint
	SortBy     ProductSortBy
	SortOrder  constants.SortOrder
}

type FindAllResult struct {
	Items []entity.Product
	Total int64
}

type ProductRepository interface {
	FindAll(context.Context, FindAllRequest) (FindAllResult, error)
	FindByID(context.Context, uint) (*entity.Product, error)
	FindByName(context.Context, string) (*entity.Product, error)
	Create(context.Context, *entity.Product) error
	Update(context.Context, *entity.Product) error
	SoftDelete(context.Context, uint) error
	Delete(context.Context, uint) error
}
