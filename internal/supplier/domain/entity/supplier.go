package entity

import (
	"time"

	userEntity "golang-fiber/internal/user/domain/entity"
	"golang-fiber/pkg/constants"
)

type Supplier struct {
	ID              uint
	Name            string
	Status          constants.SupplierStatus
	CreatedByUserID *uint
	UpdatedByUserID *uint
	CreatedByUser   *userEntity.User
	UpdatedByUser   *userEntity.User
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
