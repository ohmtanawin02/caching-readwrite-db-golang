package entity

import (
	"time"

	supplierEntity "golang-fiber/internal/supplier/domain/entity"
	userEntity "golang-fiber/internal/user/domain/entity"
	"golang-fiber/pkg/constants"
)

type Product struct {
	ID                uint
	Name              string
	Price             float64
	Stock             int
	Status            constants.ProductStatus
	SupplierID        *uint
	Supplier          *supplierEntity.Supplier
	CreatedByUserID   *uint
	UpdatedByUserID   *uint
	CreatedByUser     *userEntity.User
	UpdatedByUser     *userEntity.User
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}
