package repository

import (
	"context"
	"fmt"

	domain "golang-fiber/internal/product/domain"
	"golang-fiber/internal/product/domain/entity"
	"golang-fiber/internal/product/infrastructure/repository/models"
	"golang-fiber/pkg/common"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ProductQueryRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

type ProductQueryRepositoryCfg struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewProductQueryRepository(cfg ProductQueryRepositoryCfg) domain.ProductQueryRepository {
	return &ProductQueryRepository{db: cfg.DB, logger: cfg.Logger}
}

func (r *ProductQueryRepository) FindAll(ctx context.Context, req domain.FindAllRequest) (domain.FindAllResult, error) {
	log := common.NewRepoLogger(ctx, "ProductQueryRepository.FindAll")

	offset := (req.Page - 1) * req.Limit
	orderClause := fmt.Sprintf("%s %s", req.SortBy.Column(), req.SortOrder.SQL())

	baseQuery := r.db.WithContext(ctx).Model(&models.Product{})
	if req.SupplierID > 0 {
		baseQuery = baseQuery.Where("supplier_id = ?", req.SupplierID)
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		log.Error().Err(err).Msg("failed to count products")
		return domain.FindAllResult{}, err
	}

	var result []models.Product
	if err := baseQuery.Order(orderClause).Limit(req.Limit).Offset(offset).Find(&result).Error; err != nil {
		log.Error().Err(err).Msg("failed to find all products")
		return domain.FindAllResult{}, err
	}

	supplierMap, err := fetchSupplierMap(ctx, r.db, result)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch suppliers")
		return domain.FindAllResult{}, err
	}

	userMap, err := fetchUserMap(ctx, r.db, result)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch users")
		return domain.FindAllResult{}, err
	}

	log.Debug().Int("count", len(result)).Int64("total", total).Msg("find all products")
	return domain.FindAllResult{
		Items: toProductEntities(result, supplierMap, userMap),
		Total: total,
	}, nil
}

func (r *ProductQueryRepository) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
	log := common.NewRepoLogger(ctx, "ProductQueryRepository.FindByID")

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
