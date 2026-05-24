package repository

import (
	"golang-fiber/internal/user/domain/entity"
	"golang-fiber/internal/user/infrastructure/repository/models"
)

func toUserEntity(m models.User) entity.User {
	u := entity.User{
		ID:        m.ID,
		Username:  m.Username,
		Email:     m.Email,
		Password:  m.Password,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
		Phone:     m.Phone,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		u.DeletedAt = &t
	}
	return u
}

func toUserModel(u entity.User) models.User {
	return models.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Phone:     u.Phone,
	}
}
