package domain

import (
	"context"

	"golang-fiber/internal/product/domain/entity"
)

type ProductApplicationQuery interface {
	FindAll(context.Context, FindAllRequest) ([]entity.Product, error)
	FindByID(context.Context, uint) (*entity.Product, error)
}
