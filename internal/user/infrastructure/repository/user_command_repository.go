package repository

import (
	"context"

	domain "golang-fiber/internal/user/domain"
	"golang-fiber/internal/user/domain/entity"
	"golang-fiber/internal/user/infrastructure/repository/models"
	"golang-fiber/pkg/common"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserCommandRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

type UserCommandRepositoryCfg struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewUserCommandRepository(cfg UserCommandRepositoryCfg) domain.UserCommandRepository {
	return &UserCommandRepository{db: cfg.DB, logger: cfg.Logger}
}

func (r *UserCommandRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	log := common.NewRepoLogger(ctx, "UserCommandRepository.FindByUsername")

	var m models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	log.Debug().Str("username", username).Msg("find by username")
	u := toUserEntity(m)
	return &u, nil
}

func (r *UserCommandRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	log := common.NewRepoLogger(ctx, "UserCommandRepository.FindByEmail")

	var m models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	log.Debug().Str("email", email).Msg("find by email")
	u := toUserEntity(m)
	return &u, nil
}

func (r *UserCommandRepository) Create(ctx context.Context, u *entity.User) error {
	log := common.NewRepoLogger(ctx, "UserCommandRepository.Create")

	m := toUserModel(*u)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return err
	}
	u.ID = m.ID
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt
	return nil
}
