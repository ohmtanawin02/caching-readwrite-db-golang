package entity

import "time"

type Product struct {
	ID         uint
	Name       string
	Price      float64
	Stock      int
	SupplierID *uint
	Supplier   *Supplier
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Supplier struct {
	ID   uint
	Name string
}
