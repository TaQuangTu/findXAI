package lockdb

import (
	"context"
	"findx/internal/helpers"
	"findx/internal/system"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

type (
	RateLimiter interface {
		AllowOnce(ctx context.Context, key string, interval time.Duration) (*redis_rate.Result, error)
		Allow(ctx context.Context, key string, limit redis_rate.Limit) (*redis_rate.Result, error)
		AllowN(ctx context.Context, key string, limit redis_rate.Limit, n int) (*redis_rate.Result, error)
		Reset(ctx context.Context, key string) error
	}

	RateLimitRunner func(context.Context, func() error) error
)

type OurRateLimiter struct {
	*redis_rate.Limiter
}

func NewOurRateLimit(redisDns string) (*OurRateLimiter, error) {
	redisDb, err := helpers.ExtractRedisDB(redisDns)
	if err != nil {
		return nil, fmt.Errorf("invalid redis db")
	}
	client := redis.NewClient(&redis.Options{
		Addr: redisDns,
		DB:   redisDb,
	})
	system.RegisterRootCloser(client.Close)
	var (
		limiter = redis_rate.NewLimiter(client)
	)
	return &OurRateLimiter{
		Limiter: limiter}, nil
}

func (l *OurRateLimiter) AllowOnce(
	ctx context.Context,
	key string, interval time.Duration,
) (*redis_rate.Result, error) {
	limit := redis_rate.Limit{
		Rate:   1,
		Burst:  1,
		Period: interval,
	}
	return l.Allow(ctx, key, limit)
}

func (l *OurRateLimiter) Allow(
	ctx context.Context,
	key string, limit redis_rate.Limit,
) (*redis_rate.Result, error) {
	return l.AllowN(ctx, key, limit, 1)
}

func (l *OurRateLimiter) AllowN(
	ctx context.Context,
	key string, limit redis_rate.Limit, n int,
) (*redis_rate.Result, error) {
	result, err := l.Limiter.AllowN(ctx, key, limit, n)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (l *OurRateLimiter) Reset(ctx context.Context, key string) error {
	return l.Limiter.Reset(ctx, key)
}
