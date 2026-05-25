package repository

import (
	"context"

	"golang-fiber/internal/supplier/domain/entity"
	"golang-fiber/internal/supplier/infrastructure/repository/models"
	userEntity "golang-fiber/internal/user/domain/entity"
	"golang-fiber/pkg/constants"
)

type userRow struct {
	ID        uint
	Username  string
	Firstname string
	Lastname  string
	Email     string
	Phone     string
}

func (r *SupplierRepository) fetchUserMap(ctx context.Context, suppliers []models.Supplier) (map[uint]userRow, error) {
	idSet := make(map[uint]struct{})
	for _, p := range suppliers {
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
	if err := r.readDB.WithContext(ctx).
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

func toSupplierEntity(m models.Supplier, userMap map[uint]userRow) entity.Supplier {
	p := entity.Supplier{
		ID:              m.ID,
		Name:            m.Name,
		Status:          constants.SupplierStatus(m.Status),
		CreatedByUserID: m.CreatedByUserID,
		UpdatedByUserID: m.UpdatedByUserID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		p.DeletedAt = &t
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

func toSupplierEntities(ms []models.Supplier, userMap map[uint]userRow) []entity.Supplier {
	entities := make([]entity.Supplier, len(ms))
	for i, m := range ms {
		entities[i] = toSupplierEntity(m, userMap)
	}
	return entities
}

func toSupplierModel(p entity.Supplier) models.Supplier {
	return models.Supplier{
		ID:              p.ID,
		Name:            p.Name,
		Status:          string(p.Status),
		CreatedByUserID: p.CreatedByUserID,
		UpdatedByUserID: p.UpdatedByUserID,
	}
}
