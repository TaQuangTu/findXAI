package lockdb

import (
	"context"
	"findx/internal/liberror"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type Option func(*OurSemaphore)

func (f Option) apply(s *OurSemaphore) {
	f(s)
}

func WithRetries(retries int) Option {
	return func(s *OurSemaphore) {
		s.retires = retries
	}
}

func WithMaxSlot(maxSlot int) Option {
	return func(s *OurSemaphore) {
		s.maxSlot = maxSlot
	}
}

func WithExpiry(timeout time.Duration) Option {
	return func(s *OurSemaphore) {
		s.timeout = timeout
	}
}

func WithRetryDelay(retryDelay time.Duration) Option {
	return func(s *OurSemaphore) {
		s.retryDelay = retryDelay
	}
}

type OurSemaphore struct {
	*redis.Client

	key string

	maxSlot    int
	timeout    time.Duration
	retryDelay time.Duration
	retires    int
}

func NewSemaphore(
	redisClient *redis.Client,
	key string,
	options ...Option,
) *OurSemaphore {
	semaphore := &OurSemaphore{
		Client: redisClient,
		key:    key,
	}
	for _, opt := range options {
		opt.apply(semaphore)
	}
	return semaphore
}

func (s *OurSemaphore) AcquireSlot(ctx context.Context, key string) (err error) {
	var (
		script = `
			local current = tonumber(redis.call("get", KEYS[1]) or "0")
			if current < tonumber(ARGV[1]) then
				current = redis.call("incr", KEYS[1])
				if current == 1 then
					redis.call("expire", KEYS[1], ARGV[2])
				end
				return 1
			else
				return 0
			end
		`
	)

	for i := 0; i < s.retires; i++ {
		var (
			res any
		)
		res, err = s.Client.Eval(ctx, script, []string{key}, s.maxSlot, int(s.timeout.Seconds())).Result()
		if err != nil {
			err = liberror.WrapMessage(err, "failed to eval redis script")
			return
		}
		val, ok := res.(int64)
		if ok && val == 1 {
			return nil
		}
		time.Sleep(s.retryDelay)
	}
	err = liberror.WrapMessage(err, "failed to acquire lock").
		WithField("max_retries", s.retires)
	return
}

func (s *OurSemaphore) ReleaseSlot(ctx context.Context) (err error) {
	err = s.Client.Decr(ctx, s.key).Err()
	if err != nil {
		err = liberror.WrapMessage(err, "failed to unlock redis semaphore")
	}
	return
}
