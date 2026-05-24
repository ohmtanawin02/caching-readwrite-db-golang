package constants

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
)

func (s ProductStatus) IsValid() bool {
	return s == ProductStatusActive || s == ProductStatusInactive
}

type SupplierStatus string

const (
	SupplierStatusActive   SupplierStatus = "active"
	SupplierStatusInactive SupplierStatus = "inactive"
)

func (s SupplierStatus) IsValid() bool {
	return s == SupplierStatusActive || s == SupplierStatusInactive
}
