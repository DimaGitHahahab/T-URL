package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"storage/internal/repository"

	"github.com/redis/go-redis/v9"
)

type redisRepo struct {
	redis *redis.Client
}

func New(cl *redis.Client) repository.CacheRepository {
	return &redisRepo{redis: cl}
}

func (r *redisRepo) AddURL(ctx context.Context, short, long string) error {
	err := r.redis.Set(ctx, short, long, time.Hour*24).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}

	err = r.redis.Set(ctx, long, short, time.Hour*24).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}

	return nil
}

func (r *redisRepo) GetLongByShort(ctx context.Context, short string) (string, error) {
	result, err := r.redis.Get(ctx, short).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to get by short URL: %w", err)
	}

	return result, nil
}

func (r *redisRepo) GetShortByLong(ctx context.Context, long string) (string, error) {
	result, err := r.redis.Get(ctx, long).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to get by long URL: %w", err)
	}

	return result, nil
}

func (r *redisRepo) DeleteByShort(ctx context.Context, short string) error {
	err := r.redis.Del(ctx, short).Err()
	return fmt.Errorf("failed to delete by short URL: %w", err)
}

func (r *redisRepo) DeleteByLong(ctx context.Context, long string) error {
	err := r.redis.Del(ctx, long).Err()
	return fmt.Errorf("failed to delete by long URL: %w", err)
}

func (r *redisRepo) ShortURLExists(ctx context.Context, short string) (bool, error) {
	result, err := r.redis.Exists(ctx, short).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if short URL exists: %w", err)
	}

	return result > 0, nil
}

func (r *redisRepo) LongURLExists(ctx context.Context, long string) (bool, error) {
	result, err := r.redis.Exists(ctx, long).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if long URL exists: %w", err)
	}

	return result > 0, nil
}
