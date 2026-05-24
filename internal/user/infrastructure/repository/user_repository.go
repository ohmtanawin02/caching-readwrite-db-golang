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

type UserRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
	logger  zerolog.Logger
}

type UserRepositoryCfg struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
	Logger  zerolog.Logger
}

func NewUserRepository(cfg UserRepositoryCfg) domain.UserRepository {
	return &UserRepository{
		readDB:  cfg.ReadDB,
		writeDB: cfg.WriteDB,
		logger:  cfg.Logger,
	}
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	log := common.NewRepoLogger(ctx, "UserRepository.FindByID")

	var m models.User
	if err := r.readDB.WithContext(ctx).First(&m, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to find user by id")
		return nil, err
	}
	u := toUserEntity(m)
	return &u, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var m models.User
	if err := r.readDB.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	u := toUserEntity(m)
	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m models.User
	if err := r.readDB.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	u := toUserEntity(m)
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *entity.User) error {
	log := common.NewRepoLogger(ctx, "UserRepository.Create")

	m := toUserModel(*u)
	if err := r.writeDB.WithContext(ctx).Create(&m).Error; err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return err
	}
	u.ID = m.ID
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt
	return nil
}
