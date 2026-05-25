package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "golang-fiber/internal/supplier/domain"
	"golang-fiber/internal/supplier/domain/entity"
	"golang-fiber/pkg/auth"
	"golang-fiber/pkg/cache"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
	"golang-fiber/pkg/database"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var ErrDuplicateSupplierName = errors.New("supplier name already exists")

type SupplierCommand struct {
	repo  domain.SupplierRepository
	cache *cache.RedisCache
}

type SupplierCommandCfg struct {
	Repo  domain.SupplierRepository
	Cache *cache.RedisCache
}

func NewSupplierCommand(cfg SupplierCommandCfg) domain.SupplierApplicationCommand {
	return &SupplierCommand{
		repo:  cfg.Repo,
		cache: cfg.Cache,
	}
}

func (c *SupplierCommand) Create(ctx context.Context, input domain.CreateSupplierInput) (*entity.Supplier, error) {
	log := common.NewAppLogger(ctx, "SupplierCommand.Create")

	if _, err := c.repo.FindByName(database.WithPrimary(ctx), input.Name); err == nil {
		return nil, ErrDuplicateSupplierName
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	supplier := &entity.Supplier{
		Name:   input.Name,
		Status: constants.SupplierStatusActive,
	}

	if userID, ok := auth.GetUserID(ctx); ok {
		supplier.CreatedByUserID = &userID
		supplier.UpdatedByUserID = &userID
	}

	if err := c.repo.Create(ctx, supplier); err != nil {
		return nil, err
	}

	go c.invalidateCache(log, "suppliers:list:*")

	return c.repo.FindByID(database.WithPrimary(ctx), supplier.ID)
}

func (c *SupplierCommand) Update(ctx context.Context, id uint, input domain.UpdateSupplierInput) (*entity.Supplier, error) {
	log := common.NewAppLogger(ctx, "SupplierCommand.Update")

	supplier, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing, err := c.repo.FindByName(database.WithPrimary(ctx), input.Name); err == nil && existing.ID != id {
		return nil, ErrDuplicateSupplierName
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	supplier.Name = input.Name
	supplier.Status = input.Status

	if userID, ok := auth.GetUserID(ctx); ok {
		supplier.UpdatedByUserID = &userID
	}

	if err := c.repo.Update(ctx, supplier); err != nil {
		return nil, err
	}

	go c.invalidateCache(log,
		fmt.Sprintf("suppliers:detail:%d", id),
		"suppliers:list:*",
	)

	return c.repo.FindByID(database.WithPrimary(ctx), id)
}

func (c *SupplierCommand) SoftDelete(ctx context.Context, id uint) error {
	log := common.NewAppLogger(ctx, "SupplierCommand.SoftDelete")

	if _, err := c.repo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := c.repo.SoftDelete(ctx, id); err != nil {
		return err
	}

	go c.invalidateCache(log,
		fmt.Sprintf("suppliers:detail:%d", id),
		"suppliers:list:*",
	)
	return nil
}

func (c *SupplierCommand) Delete(ctx context.Context, id uint) error {
	log := common.NewAppLogger(ctx, "SupplierCommand.Delete")

	if _, err := c.repo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := c.repo.Delete(ctx, id); err != nil {
		return err
	}

	go c.invalidateCache(log,
		fmt.Sprintf("suppliers:detail:%d", id),
		"suppliers:list:*",
	)
	return nil
}

func (c *SupplierCommand) invalidateCache(log zerolog.Logger, patterns ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, p := range patterns {
		if err := c.cache.DelPattern(ctx, p); err != nil {
			log.Warn().Err(err).Str("pattern", p).Msg("failed to invalidate cache")
		}
	}
}
