package repository

import (
	"context"
	"fmt"

	domain "golang-fiber/internal/supplier/domain"
	"golang-fiber/internal/supplier/domain/entity"
	"golang-fiber/internal/supplier/infrastructure/repository/models"
	"golang-fiber/pkg/common"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type SupplierQueryRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

type SupplierQueryRepositoryCfg struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewSupplierQueryRepository(cfg SupplierQueryRepositoryCfg) domain.SupplierQueryRepository {
	return &SupplierQueryRepository{db: cfg.DB, logger: cfg.Logger}
}

func (r *SupplierQueryRepository) FindAll(ctx context.Context, req domain.FindAllRequest) (domain.FindAllResult, error) {
	log := common.NewRepoLogger(ctx, "SupplierQueryRepository.FindAll")

	offset := (req.Page - 1) * req.Limit
	orderClause := fmt.Sprintf("%s %s", req.SortBy.Column(), req.SortOrder.SQL())

	baseQuery := r.db.WithContext(ctx).Model(&models.Supplier{})

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		log.Error().Err(err).Msg("failed to count suppliers")
		return domain.FindAllResult{}, err
	}

	var result []models.Supplier
	if err := baseQuery.Order(orderClause).Limit(req.Limit).Offset(offset).Find(&result).Error; err != nil {
		log.Error().Err(err).Msg("failed to find all suppliers")
		return domain.FindAllResult{}, err
	}

	userMap, err := fetchUserMap(ctx, r.db, result)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch users")
		return domain.FindAllResult{}, err
	}

	log.Debug().Int("count", len(result)).Int64("total", total).Msg("find all suppliers")
	return domain.FindAllResult{
		Items: toSupplierEntities(result, userMap),
		Total: total,
	}, nil
}

func (r *SupplierQueryRepository) FindByID(ctx context.Context, id uint) (*entity.Supplier, error) {
	log := common.NewRepoLogger(ctx, "SupplierQueryRepository.FindByID")

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
