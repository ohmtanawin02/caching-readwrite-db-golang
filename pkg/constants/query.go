package constants

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func (s SortOrder) IsValid() bool {
	return s == SortOrderAsc || s == SortOrderDesc
}

func (s SortOrder) SQL() string {
	if s == SortOrderDesc {
		return "DESC"
	}
	return "ASC"
}

type ProductSortBy string

const (
	ProductSortByID        ProductSortBy = "id"
	ProductSortByName      ProductSortBy = "name"
	ProductSortByPrice     ProductSortBy = "price"
	ProductSortByStock     ProductSortBy = "stock"
	ProductSortByCreatedAt ProductSortBy = "created_at"
)

func (s ProductSortBy) IsValid() bool {
	switch s {
	case ProductSortByID, ProductSortByName, ProductSortByPrice,
		ProductSortByStock, ProductSortByCreatedAt:
		return true
	}
	return false
}

func (s ProductSortBy) Column() string {
	return string(s)
}
