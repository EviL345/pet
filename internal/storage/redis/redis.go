package redis

import (
	"blog/internal/config"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *cache.Cache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": cfg.Redis.Address,
		},
		DB:       cfg.Redis.DB,
		Password: cfg.Redis.Password,
	})

	cache := cache.New(&cache.Options{
		Redis: ring,
	})

	return cache
}
