package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID              uint           `gorm:"primaryKey"`
	Name            string         `gorm:"not null"`
	Price           float64        `gorm:"not null"`
	Stock           int
	Status          string         `gorm:"not null;default:active"`
	SupplierID      *uint          `gorm:"index"`
	CreatedByUserID *uint          `gorm:"index"`
	UpdatedByUserID *uint          `gorm:"index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}
