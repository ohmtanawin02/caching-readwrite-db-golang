package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "golang-fiber/internal/product/domain"
	"golang-fiber/internal/product/domain/entity"
	"golang-fiber/pkg/cache"
	"golang-fiber/pkg/common"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var ErrDuplicateProductName = errors.New("product name already exists")

type ProductCommand struct {
	repo  domain.ProductRepository
	cache *cache.RedisCache
}

type ProductCommandCfg struct {
	Repo  domain.ProductRepository
	Cache *cache.RedisCache
}

func NewProductCommand(cfg ProductCommandCfg) domain.ProductApplicationCommand {
	return &ProductCommand{
		repo:  cfg.Repo,
		cache: cfg.Cache,
	}
}

func (c *ProductCommand) Create(ctx context.Context, input domain.CreateProductInput) (*entity.Product, error) {
	log := common.NewAppLogger(ctx, "ProductCommand.Create")

	if _, err := c.repo.FindByName(ctx, input.Name); err == nil {
		return nil, ErrDuplicateProductName
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	product := &entity.Product{
		Name:       input.Name,
		Price:      input.Price,
		Stock:      input.Stock,
		SupplierID: input.SupplierID,
	}

	if err := c.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	go c.invalidateCache(log, "products:list:*")
	return product, nil
}

func (c *ProductCommand) Update(ctx context.Context, id uint, input domain.UpdateProductInput) (*entity.Product, error) {
	log := common.NewAppLogger(ctx, "ProductCommand.Update")

	product, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing, err := c.repo.FindByName(ctx, input.Name); err == nil && existing.ID != id {
		return nil, ErrDuplicateProductName
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	product.Name = input.Name
	product.Price = input.Price
	product.Stock = input.Stock

	if err := c.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	go c.invalidateCache(log,
		fmt.Sprintf("products:detail:%d", id),
		"products:list:*",
	)
	return product, nil
}

func (c *ProductCommand) Delete(ctx context.Context, id uint) error {
	log := common.NewAppLogger(ctx, "ProductCommand.Delete")

	if _, err := c.repo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := c.repo.Delete(ctx, id); err != nil {
		return err
	}

	go c.invalidateCache(log,
		fmt.Sprintf("products:detail:%d", id),
		"products:list:*",
	)
	return nil
}

func (c *ProductCommand) invalidateCache(log zerolog.Logger, patterns ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, p := range patterns {
		if err := c.cache.DelPattern(ctx, p); err != nil {
			log.Warn().Err(err).Str("pattern", p).Msg("failed to invalidate cache")
		}
	}
}
