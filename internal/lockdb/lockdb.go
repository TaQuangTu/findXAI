package lockdb

import (
	"time"

	redis "github.com/go-redis/redis/v8"
	redsync "github.com/go-redsync/redsync/v4"
)

type OurLockDb struct {
	*redsync.Redsync
	redisClient *redis.Client

	defaultLockTimeout time.Duration
	defaultRetryDelay  time.Duration
}

func (rs *OurLockDb) NewMutex(key string, options ...redsync.Option) *OurMutex {
	return NewMutex(rs.Redsync.NewMutex(key, options...))
}

func (rs *OurLockDb) NewSemaphore(key string, options ...Option) *OurSemaphore {
	return NewSemaphore(
		rs.redisClient,
		key,
		options...)
}
