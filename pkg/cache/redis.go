package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const DefaultTTL = 5 * time.Minute

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string, dest any) error {
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(val, dest)
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, b, ttl).Err()
}

// DelPattern ลบ key ทั้งหมดที่ match pattern เช่น "products:*"
func (r *RedisCache) DelPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
