package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"tms-core-service/internal/domain/cache"

	"github.com/redis/go-redis/v9"
)

type cacheRepo struct {
	client *redis.Client
}

// NewCacheRepository creates a new cache repository
func NewCacheRepository(client *redis.Client) cache.CacheRepository {
	return &cacheRepo{client: client}
}

// Get retrieves a value from cache
func (r *cacheRepo) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil // Key doesn't exist
		}
		return "", err
	}
	return val, nil
}

// Set stores a value in cache with expiration
func (r *cacheRepo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Convert value to JSON if it's not a string
	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		jsonData, err := json.Marshal(value)
		if err != nil {
			return err
		}
		data = string(jsonData)
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

// Delete removes a value from cache
func (r *cacheRepo) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in cache
func (r *cacheRepo) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// SetNX sets a value only if it doesn't exist (atomic)
func (r *cacheRepo) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	// Convert value to JSON if it's not a string
	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		jsonData, err := json.Marshal(value)
		if err != nil {
			return false, err
		}
		data = string(jsonData)
	}

	return r.client.SetNX(ctx, key, data, expiration).Result()
}
