package auth

import (
	"blog/internal/models"
	"context"
	"fmt"
	"github.com/go-redis/cache/v9"
	"time"
)

const basePrefix = "api-auth"

type RedisRepo struct {
	redis *cache.Cache
}

func NewRedisRepo(redis *cache.Cache) *RedisRepo {
	return &RedisRepo{redis: redis}
}

func (r *RedisRepo) GetUserById(id string) (*models.User, error) {
	user := &models.User{}
	key := generateUserKey(id)

	err := r.redis.Get(context.Background(), key, user)
	if err != nil {
		return nil, fmt.Errorf("RedisRepo.GetUserById.redis.Get: %w", err)
	}

	return user, nil
}

func (r *RedisRepo) SetUser(user *models.User, ttl int) error {
	if err := r.redis.Set(&cache.Item{
		Key:   generateUserKey(user.ID.String()),
		Value: user,
		TTL:   time.Duration(ttl) * time.Second,
	}); err != nil {
		return fmt.Errorf("RedisRepo.SetUser.redis.Set: %w", err)
	}

	return nil
}

func generateUserKey(id string) string {
	return fmt.Sprintf("%s: %s", basePrefix, id)
}
