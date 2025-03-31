package lockdb

import (
	"context"
	"findx/internal/helpers"
	"findx/internal/system"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"time"
)

type ILockDb interface {
	Lock(ctx context.Context, key string, timeout, retryDelay time.Duration) (*OurMutex, error)
	LockSimple(ctx context.Context, key string, params ...any) (*OurMutex, error)
	TryLock(ctx context.Context, key string, params ...any) (*OurMutex, error)
}

type LockDbRedis struct {
	*OurRedSync
}

func NewLockDbRedis(redisDns string) (*LockDbRedis, error) {
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
		systempool = goredis.NewPool(client)
	)
	return &LockDbRedis{
		OurRedSync: &OurRedSync{
			Redsync:            redsync.New(systempool),
			defaultLockTimeout: 0,
			defaultRetryDelay:  0,
		},
	}, nil
}

func (l *LockDbRedis) Lock(ctx context.Context, key string, timeout, retryDelay time.Duration) (*OurMutex, error) {
	var (
		mux = l.NewMutex(
			key,
			redsync.WithTries(int(timeout/retryDelay)),
			redsync.WithExpiry(timeout),
			redsync.WithRetryDelay(retryDelay),
		)
	)
	if err := mux.LockContext(ctx); err != nil {
		return nil, err
	}
	return mux, nil
}

func (l *LockDbRedis) LockSimple(ctx context.Context, key string, params ...any) (*OurMutex, error) {
	if len(params) > 0 {
		key = fmt.Sprintf(key, params...)
	}
	return l.Lock(
		ctx,
		key,
		l.defaultLockTimeout,
		l.defaultRetryDelay,
	)
}

func (l *LockDbRedis) TryLock(ctx context.Context, key string, params ...any) (*OurMutex, error) {
	if len(params) > 0 {
		key = fmt.Sprintf(key, params...)
	}
	var (
		mux = l.NewMutex(key)
	)
	if err := mux.TryLockContext(ctx); err != nil {
		return nil, err
	}
	return mux, nil
}
