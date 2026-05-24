package domain

import (
	"context"

	"golang-fiber/internal/user/domain/entity"
)

type UserRepository interface {
	FindByID(context.Context, uint) (*entity.User, error)
	FindByUsername(context.Context, string) (*entity.User, error)
	FindByEmail(context.Context, string) (*entity.User, error)
	Create(context.Context, *entity.User) error
}
