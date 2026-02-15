package redis

import (
	"context"
	"fmt"
	"log"

	"tms-core-service/internal/config"

	"github.com/redis/go-redis/v9"
)

// NewConnection creates a new Redis connection
func NewConnection(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// Ping to verify connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("✅ Redis connection established")

	return client, nil
}
