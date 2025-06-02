package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
)

var ErrNotFound = errors.New("not found")
const defaultTTL = time.Hour

type RedisCache struct {
	client *redis.Client
	ttl time.Duration
}

func NewRedisCache(ctx context.Context) (*RedisCache, error) {
	options, err := redis.ParseURL(config.Config.RedisConn)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(options)
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisCache{
		client: rdb,
		ttl: defaultTTL,
	}, nil
}

func (r *RedisCache) Close() {
	if err :=  r.client.Close(); err != nil {
		logger.Logger.Error("failed to close redis connection",
			slog.String("error", err.Error()),
			slog.String("where", "RedisCache.Close"))
	}
}

func (r *RedisCache) Set(ctx context.Context, model, id string, value any) error {
	key := fmt.Sprintf("%s:%s", model, id)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, data, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}
	return nil
}

func (r *RedisCache) Get(ctx context.Context, model, id string, dst any) error {
	key := fmt.Sprintf("%s:%s", model, id)
	var data string
	err := r.client.Get(ctx, key).Scan(&data)
	if errors.Is(err, redis.Nil) {
		return ErrNotFound
	} else if err != nil {
		return fmt.Errorf("redis.Get: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dst); err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, model, id string) error {
	key := fmt.Sprintf("%s:%s", model, id)
	_, err := r.client.Del(ctx, key).Result()
	return err
}

func (r *RedisCache) Clear(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}
