package queries

import (
	"context"
	"errors"
	"fmt"

	domain "golang-fiber/internal/product/domain"
	"golang-fiber/internal/product/domain/entity"
	"golang-fiber/pkg/cache"
	"golang-fiber/pkg/common"

	"github.com/redis/go-redis/v9"
)

type ProductQuery struct {
	repo  domain.ProductRepository
	cache *cache.RedisCache
}

type ProductQueryCfg struct {
	Repo  domain.ProductRepository
	Cache *cache.RedisCache
}

func NewProductQuery(cfg ProductQueryCfg) domain.ProductApplicationQuery {
	return &ProductQuery{
		repo:  cfg.Repo,
		cache: cfg.Cache,
	}
}

func cacheKeyList(req domain.FindAllRequest) string {
	return fmt.Sprintf("products:list:supplier:%d:page:%d:limit:%d:sort:%s:%s",
		req.SupplierID, req.Page, req.Limit, req.SortBy, req.SortOrder)
}

func (q *ProductQuery) FindAll(ctx context.Context, req domain.FindAllRequest) ([]entity.Product, error) {
	log := common.NewAppLogger(ctx, "ProductQuery.FindAll")
	key := cacheKeyList(req)

	var cached []entity.Product
	if err := q.cache.Get(ctx, key, &cached); err == nil {
		log.Debug().Str("key", key).Msg("cache hit")
		return cached, nil
	} else if !errors.Is(err, redis.Nil) {
		log.Warn().Err(err).Msg("redis error, fallthrough to db")
	}

	log.Debug().Str("key", key).Msg("cache miss")
	products, err := q.repo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := q.cache.Set(context.Background(), key, products, cache.DefaultTTL); err != nil {
			log.Warn().Err(err).Msg("failed to set cache")
		}
	}()

	return products, nil
}

func (q *ProductQuery) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
	log := common.NewAppLogger(ctx, "ProductQuery.FindByID")
	key := fmt.Sprintf("products:detail:%d", id)

	var cached entity.Product
	if err := q.cache.Get(ctx, key, &cached); err == nil {
		log.Debug().Str("key", key).Msg("cache hit")
		return &cached, nil
	}

	product, err := q.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := q.cache.Set(context.Background(), key, product, cache.DefaultTTL); err != nil {
			log.Warn().Err(err).Msg("failed to set cache")
		}
	}()

	return product, nil
}
