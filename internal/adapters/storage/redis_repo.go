package storage

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisTokenRepo struct {
	client *redis.Client
}

func NewRedisTokenRepo(client *redis.Client) *RedisTokenRepo {
	return &RedisTokenRepo{client: client}
}

func (r *RedisTokenRepo) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, tokenID)
	return r.client.Set(ctx, key, "valid", expiresIn).Err()
}

func (r *RedisTokenRepo) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, tokenID)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisTokenRepo) ValidateRefreshToken(ctx context.Context, userID string, tokenID string) (bool, error) {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, tokenID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "valid", nil
}
