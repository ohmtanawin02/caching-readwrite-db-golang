package domain

type SupplierSortBy string

const (
	SupplierSortByID        SupplierSortBy = "id"
	SupplierSortByName      SupplierSortBy = "name"
	SupplierSortByCreatedAt SupplierSortBy = "created_at"
)

func (s SupplierSortBy) IsValid() bool {
	switch s {
	case SupplierSortByID, SupplierSortByName, SupplierSortByCreatedAt:
		return true
	}
	return false
}

func (s SupplierSortBy) Column() string {
	return string(s)
}
