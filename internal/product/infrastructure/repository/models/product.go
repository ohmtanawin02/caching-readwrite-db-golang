package models

import "time"

type Product struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string    `gorm:"not null"`
	Price      float64   `gorm:"not null"`
	Stock      int
	SupplierID *uint     `gorm:"index"`
	Supplier   *Supplier `gorm:"foreignKey:SupplierID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Product) TableName() string {
	return "products"
}
