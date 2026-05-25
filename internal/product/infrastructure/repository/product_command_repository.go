package repository

import (
	"context"

	domain "golang-fiber/internal/product/domain"
	"golang-fiber/internal/product/domain/entity"
	"golang-fiber/internal/product/infrastructure/repository/models"
	"golang-fiber/pkg/common"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ProductCommandRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

type ProductCommandRepositoryCfg struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewProductCommandRepository(cfg ProductCommandRepositoryCfg) domain.ProductCommandRepository {
	return &ProductCommandRepository{db: cfg.DB, logger: cfg.Logger}
}

func (r *ProductCommandRepository) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
	log := common.NewRepoLogger(ctx, "ProductCommandRepository.FindByID")

	var result models.Product
	if err := r.db.WithContext(ctx).First(&result, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to find product")
		return nil, err
	}

	supplierMap, err := fetchSupplierMap(ctx, r.db, []models.Product{result})
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch supplier")
		return nil, err
	}

	userMap, err := fetchUserMap(ctx, r.db, []models.Product{result})
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch users")
		return nil, err
	}

	p := toProductEntity(result, supplierMap, userMap)
	return &p, nil
}

func (r *ProductCommandRepository) FindByName(ctx context.Context, name string) (*entity.Product, error) {
	log := common.NewRepoLogger(ctx, "ProductCommandRepository.FindByName")

	var result models.Product
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&result).Error; err != nil {
		return nil, err
	}

	log.Debug().Str("name", name).Msg("find by name")
	p := toProductEntity(result, nil, nil)
	return &p, nil
}

func (r *ProductCommandRepository) Create(ctx context.Context, p *entity.Product) error {
	log := common.NewRepoLogger(ctx, "ProductCommandRepository.Create")

	m := toProductModel(*p)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		log.Error().Err(err).Msg("failed to create product")
		return err
	}

	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *ProductCommandRepository) Update(ctx context.Context, p *entity.Product) error {
	log := common.NewRepoLogger(ctx, "ProductCommandRepository.Update")

	m := toProductModel(*p)
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		log.Error().Err(err).Uint("id", p.ID).Msg("failed to update product")
		return err
	}
	return nil
}

func (r *ProductCommandRepository) SoftDelete(ctx context.Context, id uint) error {
	log := common.NewRepoLogger(ctx, "ProductCommandRepository.SoftDelete")

	if err := r.db.WithContext(ctx).Delete(&models.Product{}, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to soft delete product")
		return err
	}
	return nil
}

func (r *ProductCommandRepository) Delete(ctx context.Context, id uint) error {
	log := common.NewRepoLogger(ctx, "ProductCommandRepository.Delete")

	if err := r.db.WithContext(ctx).Unscoped().Delete(&models.Product{}, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to hard delete product")
		return err
	}
	return nil
}
