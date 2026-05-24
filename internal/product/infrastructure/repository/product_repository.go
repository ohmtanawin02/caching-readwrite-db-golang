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

func (r *ProductRepository) FindAll(ctx context.Context, req domain.FindAllRequest) ([]entity.Product, error) {
	log := common.NewRepoLogger(ctx, "ProductRepository.FindAll")

	offset := (req.Page - 1) * req.Limit
	orderClause := fmt.Sprintf("%s %s", req.SortBy.Column(), req.SortOrder.SQL())

	query := r.readDB.WithContext(ctx).
		Preload("Supplier").
		Order(orderClause).
		Limit(req.Limit).
		Offset(offset)

	if req.SupplierID > 0 {
		query = query.Where("supplier_id = ?", req.SupplierID)
	}

	var result []models.Product
	if err := query.Find(&result).Error; err != nil {
		log.Error().Err(err).Msg("failed to find all products")
		return nil, err
	}

	log.Debug().Int("count", len(result)).Msg("find all products")
	return toProductEntities(result), nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
	log := common.NewRepoLogger(ctx, "ProductRepository.FindByID")

	var result models.Product
	if err := r.readDB.WithContext(ctx).Preload("Supplier").First(&result, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to find product")
		return nil, err
	}

	p := toProductEntity(result)
	return &p, nil
}

func (r *ProductRepository) FindByName(ctx context.Context, name string) (*entity.Product, error) {
	log := common.NewRepoLogger(ctx, "ProductRepository.FindByName")

	var result models.Product
	if err := r.readDB.WithContext(ctx).Where("name = ?", name).First(&result).Error; err != nil {
		return nil, err
	}

	log.Debug().Str("name", name).Msg("find by name")
	p := toProductEntity(result)
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

func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
	log := common.NewRepoLogger(ctx, "ProductRepository.Delete")

	if err := r.writeDB.WithContext(ctx).Delete(&models.Product{}, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to delete product")
		return err
	}
	return nil
}

func toProductEntity(m models.Product) entity.Product {
	p := entity.Product{
		ID:         m.ID,
		Name:       m.Name,
		Price:      m.Price,
		Stock:      m.Stock,
		SupplierID: m.SupplierID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
	if m.Supplier != nil {
		p.Supplier = &entity.Supplier{ID: m.Supplier.ID, Name: m.Supplier.Name}
	}
	return p
}

func toProductEntities(ms []models.Product) []entity.Product {
	entities := make([]entity.Product, len(ms))
	for i, m := range ms {
		entities[i] = toProductEntity(m)
	}
	return entities
}

func toProductModel(p entity.Product) models.Product {
	return models.Product{
		ID:         p.ID,
		Name:       p.Name,
		Price:      p.Price,
		Stock:      p.Stock,
		SupplierID: p.SupplierID,
	}
}
