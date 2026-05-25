package repository

import (
	"context"

	"golang-fiber/internal/product/domain/entity"
	"golang-fiber/internal/product/infrastructure/repository/models"
	supplierEntity "golang-fiber/internal/supplier/domain/entity"
	userEntity "golang-fiber/internal/user/domain/entity"
	"golang-fiber/pkg/constants"

	"gorm.io/gorm"
)

// supplierRow — local scan struct (ไม่ import supplier infrastructure)
type supplierRow struct {
	ID     uint
	Name   string
	Status string
}

// userRow — local scan struct (ไม่ import user infrastructure)
type userRow struct {
	ID        uint
	Username  string
	Firstname string
	Lastname  string
	Email     string
	Phone     string
}

func fetchSupplierMap(ctx context.Context, db *gorm.DB, products []models.Product) (map[uint]supplierRow, error) {
	idSet := make(map[uint]struct{})
	for _, p := range products {
		if p.SupplierID != nil {
			idSet[*p.SupplierID] = struct{}{}
		}
	}
	if len(idSet) == 0 {
		return nil, nil
	}

	ids := make([]uint, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	var rows []supplierRow
	if err := db.WithContext(ctx).
		Table("suppliers").
		Select("id, name, status").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	m := make(map[uint]supplierRow, len(rows))
	for _, s := range rows {
		m[s.ID] = s
	}
	return m, nil
}

func fetchUserMap(ctx context.Context, db *gorm.DB, products []models.Product) (map[uint]userRow, error) {
	idSet := make(map[uint]struct{})
	for _, p := range products {
		if p.CreatedByUserID != nil {
			idSet[*p.CreatedByUserID] = struct{}{}
		}
		if p.UpdatedByUserID != nil {
			idSet[*p.UpdatedByUserID] = struct{}{}
		}
	}
	if len(idSet) == 0 {
		return nil, nil
	}

	ids := make([]uint, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	var rows []userRow
	if err := db.WithContext(ctx).
		Table("users").
		Select("id, username, firstname, lastname, email, phone").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	m := make(map[uint]userRow, len(rows))
	for _, u := range rows {
		m[u.ID] = u
	}
	return m, nil
}

func toProductEntity(m models.Product, supplierMap map[uint]supplierRow, userMap map[uint]userRow) entity.Product {
	p := entity.Product{
		ID:              m.ID,
		Name:            m.Name,
		Price:           m.Price,
		Stock:           m.Stock,
		Status:          constants.ProductStatus(m.Status),
		SupplierID:      m.SupplierID,
		CreatedByUserID: m.CreatedByUserID,
		UpdatedByUserID: m.UpdatedByUserID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		p.DeletedAt = &t
	}
	if m.SupplierID != nil {
		if s, ok := supplierMap[*m.SupplierID]; ok {
			p.Supplier = &supplierEntity.Supplier{
				ID:     s.ID,
				Name:   s.Name,
				Status: constants.SupplierStatus(s.Status),
			}
		}
	}
	if m.CreatedByUserID != nil {
		if u, ok := userMap[*m.CreatedByUserID]; ok {
			p.CreatedByUser = &userEntity.User{
				ID: u.ID, Username: u.Username,
				Firstname: u.Firstname, Lastname: u.Lastname,
				Email: u.Email, Phone: u.Phone,
			}
		}
	}
	if m.UpdatedByUserID != nil {
		if u, ok := userMap[*m.UpdatedByUserID]; ok {
			p.UpdatedByUser = &userEntity.User{
				ID: u.ID, Username: u.Username,
				Firstname: u.Firstname, Lastname: u.Lastname,
				Email: u.Email, Phone: u.Phone,
			}
		}
	}
	return p
}

func toProductEntities(ms []models.Product, supplierMap map[uint]supplierRow, userMap map[uint]userRow) []entity.Product {
	entities := make([]entity.Product, len(ms))
	for i, m := range ms {
		entities[i] = toProductEntity(m, supplierMap, userMap)
	}
	return entities
}

func toProductModel(p entity.Product) models.Product {
	return models.Product{
		ID:              p.ID,
		Name:            p.Name,
		Price:           p.Price,
		Stock:           p.Stock,
		Status:          string(p.Status),
		SupplierID:      p.SupplierID,
		CreatedByUserID: p.CreatedByUserID,
		UpdatedByUserID: p.UpdatedByUserID,
	}
}
