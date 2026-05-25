package repository

import (
	"context"

	domain "golang-fiber/internal/supplier/domain"
	"golang-fiber/internal/supplier/domain/entity"
	"golang-fiber/internal/supplier/infrastructure/repository/models"
	"golang-fiber/pkg/common"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type SupplierCommandRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

type SupplierCommandRepositoryCfg struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewSupplierCommandRepository(cfg SupplierCommandRepositoryCfg) domain.SupplierCommandRepository {
	return &SupplierCommandRepository{db: cfg.DB, logger: cfg.Logger}
}

func (r *SupplierCommandRepository) FindByID(ctx context.Context, id uint) (*entity.Supplier, error) {
	log := common.NewRepoLogger(ctx, "SupplierCommandRepository.FindByID")

	var result models.Supplier
	if err := r.db.WithContext(ctx).First(&result, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to find supplier")
		return nil, err
	}

	userMap, err := fetchUserMap(ctx, r.db, []models.Supplier{result})
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch users")
		return nil, err
	}

	p := toSupplierEntity(result, userMap)
	return &p, nil
}

func (r *SupplierCommandRepository) FindByName(ctx context.Context, name string) (*entity.Supplier, error) {
	log := common.NewRepoLogger(ctx, "SupplierCommandRepository.FindByName")

	var result models.Supplier
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&result).Error; err != nil {
		return nil, err
	}

	log.Debug().Str("name", name).Msg("find by name")
	p := toSupplierEntity(result, nil)
	return &p, nil
}

func (r *SupplierCommandRepository) Create(ctx context.Context, p *entity.Supplier) error {
	log := common.NewRepoLogger(ctx, "SupplierCommandRepository.Create")

	m := toSupplierModel(*p)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		log.Error().Err(err).Msg("failed to create supplier")
		return err
	}

	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *SupplierCommandRepository) Update(ctx context.Context, p *entity.Supplier) error {
	log := common.NewRepoLogger(ctx, "SupplierCommandRepository.Update")

	m := toSupplierModel(*p)
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		log.Error().Err(err).Uint("id", p.ID).Msg("failed to update supplier")
		return err
	}
	return nil
}

func (r *SupplierCommandRepository) SoftDelete(ctx context.Context, id uint) error {
	log := common.NewRepoLogger(ctx, "SupplierCommandRepository.SoftDelete")

	if err := r.db.WithContext(ctx).Delete(&models.Supplier{}, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to soft delete supplier")
		return err
	}
	return nil
}

func (r *SupplierCommandRepository) Delete(ctx context.Context, id uint) error {
	log := common.NewRepoLogger(ctx, "SupplierCommandRepository.Delete")

	if err := r.db.WithContext(ctx).Unscoped().Delete(&models.Supplier{}, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to hard delete supplier")
		return err
	}
	return nil
}
