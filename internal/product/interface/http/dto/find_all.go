package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	domain "golang-fiber/internal/product/domain"
	"golang-fiber/pkg/constants"
)

type FindAllRequest struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	SupplierID uint   `query:"supplier_id"`
	SortBy     string `query:"sort_by"`
	SortOrder  string `query:"sort_order"`
}

func ValidateFindAllRequest(c *fiber.Ctx, v *validator.Validate) (*FindAllRequest, error) {
	req := new(FindAllRequest)
	if err := c.QueryParser(req); err != nil {
		return nil, err
	}
	return req, nil
}

func ToFindAllDomainRequest(req *FindAllRequest) domain.FindAllRequest {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	sortBy := domain.ProductSortBy(req.SortBy)
	if !sortBy.IsValid() {
		sortBy = domain.ProductSortByID
	}

	sortOrder := constants.SortOrder(req.SortOrder)
	if !sortOrder.IsValid() {
		sortOrder = constants.SortOrderAsc
	}

	return domain.FindAllRequest{
		Page:       page,
		Limit:      limit,
		SupplierID: req.SupplierID,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}
}
