package dto

import (
	"math"
	"time"

	"golang-fiber/internal/supplier/domain"
	"golang-fiber/internal/supplier/domain/entity"
	"golang-fiber/pkg/constants"
)

type UserBriefResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type SupplierResponse struct {
	ID        uint                     `json:"id"`
	Name      string                   `json:"name"`
	Status    constants.SupplierStatus `json:"status"`
	CreatedBy *UserBriefResponse       `json:"created_by,omitempty"`
	UpdatedBy *UserBriefResponse       `json:"updated_by,omitempty"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
	DeletedAt *time.Time               `json:"deleted_at,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type SupplierListResponse struct {
	Items []SupplierResponse `json:"items"`
	Meta  PaginationMeta     `json:"meta"`
}

func ToSupplierResponse(p entity.Supplier) SupplierResponse {
	r := SupplierResponse{
		ID:        p.ID,
		Name:      p.Name,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,
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

func ToSupplierListResponse(result domain.FindAllResult, req domain.FindAllRequest) SupplierListResponse {
	items := make([]SupplierResponse, len(result.Items))
	for i, p := range result.Items {
		items[i] = ToSupplierResponse(p)
	}

	totalPages := 0
	if req.Limit > 0 {
		totalPages = int(math.Ceil(float64(result.Total) / float64(req.Limit)))
	}

	return SupplierListResponse{
		Items: items,
		Meta: PaginationMeta{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      result.Total,
			TotalPages: totalPages,
		},
	}
}
