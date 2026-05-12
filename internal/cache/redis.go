package cache

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/NickSFU/shortlink-service/internal/config"
)

func NewRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	return client
}

func Ping(client *redis.Client) error {
	return client.Ping(context.Background()).Err()
}
