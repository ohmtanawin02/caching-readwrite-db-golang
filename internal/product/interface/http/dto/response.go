package dto

import (
	"math"
	"time"

	"golang-fiber/internal/product/domain"
	"golang-fiber/internal/product/domain/entity"
	"golang-fiber/pkg/constants"
)

type UserBriefResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type SupplierResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProductResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Price       float64                 `json:"price"`
	Stock       int                     `json:"stock"`
	Status      constants.ProductStatus `json:"status"`
	SupplierID  *uint                   `json:"supplier_id"`
	Supplier    *SupplierResponse       `json:"supplier,omitempty"`
	CreatedBy   *UserBriefResponse      `json:"created_by,omitempty"`
	UpdatedBy   *UserBriefResponse      `json:"updated_by,omitempty"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
	DeletedAt   *time.Time              `json:"deleted_at,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type ProductListResponse struct {
	Items []ProductResponse `json:"items"`
	Meta  PaginationMeta    `json:"meta"`
}

func ToProductResponse(p entity.Product) ProductResponse {
	r := ProductResponse{
		ID:         p.ID,
		Name:       p.Name,
		Price:      p.Price,
		Stock:      p.Stock,
		Status:     p.Status,
		SupplierID: p.SupplierID,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		DeletedAt:  p.DeletedAt,
	}
	if p.Supplier != nil {
		r.Supplier = &SupplierResponse{ID: p.Supplier.ID, Name: p.Supplier.Name}
	}
	if p.CreatedByUser != nil {
		r.CreatedBy = &UserBriefResponse{
			ID: p.CreatedByUser.ID, Username: p.CreatedByUser.Username,
			Firstname: p.CreatedByUser.Firstname, Lastname: p.CreatedByUser.Lastname,
		}
	}
	if p.UpdatedByUser != nil {
		r.UpdatedBy = &UserBriefResponse{
			ID: p.UpdatedByUser.ID, Username: p.UpdatedByUser.Username,
			Firstname: p.UpdatedByUser.Firstname, Lastname: p.UpdatedByUser.Lastname,
		}
	}
	return r
}

func ToProductListResponse(result domain.FindAllResult, req domain.FindAllRequest) ProductListResponse {
	items := make([]ProductResponse, len(result.Items))
	for i, p := range result.Items {
		items[i] = ToProductResponse(p)
	}

	totalPages := 0
	if req.Limit > 0 {
		totalPages = int(math.Ceil(float64(result.Total) / float64(req.Limit)))
	}

	return ProductListResponse{
		Items: items,
		Meta: PaginationMeta{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      result.Total,
			TotalPages: totalPages,
		},
	}
}
