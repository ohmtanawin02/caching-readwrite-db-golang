package repository

import (
	"context"
	"fmt"

	domain "golang-fiber/internal/supplier/domain"
	"golang-fiber/internal/supplier/domain/entity"
	"golang-fiber/internal/supplier/infrastructure/repository/models"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/database"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type SupplierRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
	logger  zerolog.Logger
}

type SupplierRepositoryCfg struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
	Logger  zerolog.Logger
}

func NewSupplierRepository(cfg SupplierRepositoryCfg) domain.SupplierRepository {
	return &SupplierRepository{
		readDB:  cfg.ReadDB,
		writeDB: cfg.WriteDB,
		logger:  cfg.Logger,
	}
}

func (r *SupplierRepository) FindAll(ctx context.Context, req domain.FindAllRequest) (domain.FindAllResult, error) {
	log := common.NewRepoLogger(ctx, "SupplierRepository.FindAll")

	offset := (req.Page - 1) * req.Limit
	orderClause := fmt.Sprintf("%s %s", req.SortBy.Column(), req.SortOrder.SQL())

	baseQuery := r.readDB.WithContext(ctx).Model(&models.Supplier{})

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

	userMap, err := r.fetchUserMap(ctx, result)
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

func (r *SupplierRepository) FindByID(ctx context.Context, id uint) (*entity.Supplier, error) {
	log := common.NewRepoLogger(ctx, "SupplierRepository.FindByID")

	db := r.readDB
	if database.UsePrimary(ctx) {
		db = r.writeDB
	}

	var result models.Supplier
	if err := db.WithContext(ctx).First(&result, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to find supplier")
		return nil, err
	}

	userMap, err := r.fetchUserMap(ctx, []models.Supplier{result})
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch users")
		return nil, err
	}

	p := toSupplierEntity(result, userMap)
	return &p, nil
}

func (r *SupplierRepository) FindByName(ctx context.Context, name string) (*entity.Supplier, error) {
	log := common.NewRepoLogger(ctx, "SupplierRepository.FindByName")

	db := r.readDB
	if database.UsePrimary(ctx) {
		db = r.writeDB
	}

	var result models.Supplier
	if err := db.WithContext(ctx).Where("name = ?", name).First(&result).Error; err != nil {
		return nil, err
	}

	log.Debug().Str("name", name).Msg("find by name")
	p := toSupplierEntity(result, nil)
	return &p, nil
}

func (r *SupplierRepository) Create(ctx context.Context, p *entity.Supplier) error {
	log := common.NewRepoLogger(ctx, "SupplierRepository.Create")

	m := toSupplierModel(*p)
	if err := r.writeDB.WithContext(ctx).Create(&m).Error; err != nil {
		log.Error().Err(err).Msg("failed to create supplier")
		return err
	}

	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *SupplierRepository) Update(ctx context.Context, p *entity.Supplier) error {
	log := common.NewRepoLogger(ctx, "SupplierRepository.Update")

	m := toSupplierModel(*p)
	if err := r.writeDB.WithContext(ctx).Save(&m).Error; err != nil {
		log.Error().Err(err).Uint("id", p.ID).Msg("failed to update supplier")
		return err
	}
	return nil
}

func (r *SupplierRepository) SoftDelete(ctx context.Context, id uint) error {
	log := common.NewRepoLogger(ctx, "SupplierRepository.SoftDelete")

	if err := r.writeDB.WithContext(ctx).Delete(&models.Supplier{}, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to soft delete supplier")
		return err
	}
	return nil
}

func (r *SupplierRepository) Delete(ctx context.Context, id uint) error {
	log := common.NewRepoLogger(ctx, "SupplierRepository.Delete")

	if err := r.writeDB.WithContext(ctx).Unscoped().Delete(&models.Supplier{}, id).Error; err != nil {
		log.Error().Err(err).Uint("id", id).Msg("failed to hard delete supplier")
		return err
	}
	return nil
}
