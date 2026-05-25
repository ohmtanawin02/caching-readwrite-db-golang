package queries

import (
	"context"
	"errors"
	"fmt"

	domain "golang-fiber/internal/supplier/domain"
	"golang-fiber/internal/supplier/domain/entity"
	"golang-fiber/pkg/cache"
	"golang-fiber/pkg/common"

	"github.com/redis/go-redis/v9"
)

type SupplierQuery struct {
	repo  domain.SupplierQueryRepository
	cache *cache.RedisCache
}

type SupplierQueryCfg struct {
	Repo  domain.SupplierQueryRepository
	Cache *cache.RedisCache
}

func NewSupplierQuery(cfg SupplierQueryCfg) domain.SupplierApplicationQuery {
	return &SupplierQuery{
		repo:  cfg.Repo,
		cache: cfg.Cache,
	}
}

func cacheKeyList(req domain.FindAllRequest) string {
	return fmt.Sprintf("suppliers:list:page:%d:limit:%d:sort:%s:%s",
		req.Page, req.Limit, req.SortBy, req.SortOrder)
}

func (q *SupplierQuery) FindAll(ctx context.Context, req domain.FindAllRequest) (domain.FindAllResult, error) {
	log := common.NewAppLogger(ctx, "SupplierQuery.FindAll")
	key := cacheKeyList(req)

	var cached domain.FindAllResult
	if err := q.cache.Get(ctx, key, &cached); err == nil {
		log.Debug().Str("key", key).Msg("cache hit")
		return cached, nil
	} else if !errors.Is(err, redis.Nil) {
		log.Warn().Err(err).Msg("redis error, fallthrough to db")
	}

	log.Debug().Str("key", key).Msg("cache miss")
	result, err := q.repo.FindAll(ctx, req)
	if err != nil {
		return domain.FindAllResult{}, err
	}

	go func() {
		if err := q.cache.Set(context.Background(), key, result, cache.DefaultTTL); err != nil {
			log.Warn().Err(err).Msg("failed to set cache")
		}
	}()

	return result, nil
}

func (q *SupplierQuery) FindByID(ctx context.Context, id uint) (*entity.Supplier, error) {
	log := common.NewAppLogger(ctx, "SupplierQuery.FindByID")
	key := fmt.Sprintf("suppliers:detail:%d", id)

	var cached entity.Supplier
	if err := q.cache.Get(ctx, key, &cached); err == nil {
		log.Debug().Str("key", key).Msg("cache hit")
		return &cached, nil
	} else if !errors.Is(err, redis.Nil) {
		log.Warn().Err(err).Msg("redis error, fallthrough to db")
	}

	supplier, err := q.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	snapshot := *supplier
	go func() {
		if err := q.cache.Set(context.Background(), key, snapshot, cache.DefaultTTL); err != nil {
			log.Warn().Err(err).Msg("failed to set cache")
		}
	}()

	return supplier, nil
}
