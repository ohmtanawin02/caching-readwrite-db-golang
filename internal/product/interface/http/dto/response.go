package dto

import (
	"time"

	"golang-fiber/internal/product/domain/entity"
)

type ProductResponse struct {
	ID         uint              `json:"id"`
	Name       string            `json:"name"`
	Price      float64           `json:"price"`
	Stock      int               `json:"stock"`
	SupplierID *uint             `json:"supplier_id"`
	Supplier   *SupplierResponse `json:"supplier,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type SupplierResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ToProductResponse(p entity.Product) ProductResponse {
	r := ProductResponse{
		ID:         p.ID,
		Name:       p.Name,
		Price:      p.Price,
		Stock:      p.Stock,
		SupplierID: p.SupplierID,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
	if p.Supplier != nil {
		r.Supplier = &SupplierResponse{ID: p.Supplier.ID, Name: p.Supplier.Name}
	}
	return r
}

func ToProductResponses(products []entity.Product) []ProductResponse {
	resp := make([]ProductResponse, len(products))
	for i, p := range products {
		resp[i] = ToProductResponse(p)
	}
	return resp
}
