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

type ProductRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
	logger  zerolog.Logger
}

type ProductRepositoryCfg struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
	Logger  zerolog.Logger
}

func NewProductRepository(cfg ProductRepositoryCfg) domain.ProductRepository {
	return &ProductRepository{
		readDB:  cfg.ReadDB,
		writeDB: cfg.WriteDB,
		logger:  cfg.Logger,
	}
}

func (r *ProductRepository) FindAll(ctx context.Context, req domain.FindAllRequest) (domain.FindAllResult, error) {
	log := common.NewRepoLogger(ctx, "ProductRepository.FindAll")

	offset := (req.Page - 1) * req.Limit
	orderClause := fmt.Sprintf("%s %s", req.SortBy.Column(), req.SortOrder.SQL())

	baseQuery := r.readDB.WithContext(ctx).Model(&models.Product{})
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

	supplierMap, err := r.fetchSupplierMap(ctx, result)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch suppliers")
		return domain.FindAllResult{}, err
	}

	userMap, err := r.fetchUserMap(ctx, result)
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

func (r *ProductRepository) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
	log := common.NewRepoLogger(ctx, "ProductRepository.FindByID")

	var result models.Product
	if err := r.readDB.WithContext(ctx).First(&result, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to find product")
		return nil, err
	}

	supplierMap, err := r.fetchSupplierMap(ctx, []models.Product{result})
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch supplier")
		return nil, err
	}

	userMap, err := r.fetchUserMap(ctx, []models.Product{result})
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch users")
		return nil, err
	}

	p := toProductEntity(result, supplierMap, userMap)
	return &p, nil
}

func (r *ProductRepository) FindByName(ctx context.Context, name string) (*entity.Product, error) {
	log := common.NewRepoLogger(ctx, "ProductRepository.FindByName")

	var result models.Product
	if err := r.readDB.WithContext(ctx).Where("name = ?", name).First(&result).Error; err != nil {
		return nil, err
	}

	log.Debug().Str("name", name).Msg("find by name")
	p := toProductEntity(result, nil, nil)
	return &p, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *entity.Product) error {
	log := common.NewRepoLogger(ctx, "ProductRepository.Create")

	m := toProductModel(*p)
	if err := r.writeDB.WithContext(ctx).Create(&m).Error; err != nil {
		log.Error().Err(err).Msg("failed to create product")
		return err
	}

	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *ProductRepository) Update(ctx context.Context, p *entity.Product) error {
	log := common.NewRepoLogger(ctx, "ProductRepository.Update")

	m := toProductModel(*p)
	if err := r.writeDB.WithContext(ctx).Save(&m).Error; err != nil {
		log.Error().Err(err).Uint("id", p.ID).Msg("failed to update product")
		return err
	}
	return nil
}

func (r *ProductRepository) SoftDelete(ctx context.Context, id uint) error {
	log := common.NewRepoLogger(ctx, "ProductRepository.SoftDelete")

	if err := r.writeDB.WithContext(ctx).Delete(&models.Product{}, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to soft delete product")
		return err
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
	log := common.NewRepoLogger(ctx, "ProductRepository.Delete")

	if err := r.writeDB.WithContext(ctx).Unscoped().Delete(&models.Product{}, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to hard delete product")
		return err
	}
	return nil
}

func (r *ProductRepository) Transact(ctx context.Context, fn func(domain.ProductRepository) error) error {
	return r.writeDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &ProductRepository{
			readDB:  r.readDB,
			writeDB: tx,
			logger:  r.logger,
		}
		return fn(txRepo)
	})
}
