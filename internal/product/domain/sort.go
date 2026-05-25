package domain

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
