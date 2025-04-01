package lockdb

import (
	"time"

	"github.com/go-redsync/redsync/v4"
)

type OurRedSync struct {
	*redsync.Redsync

	defaultLockTimeout time.Duration
	defaultRetryDelay  time.Duration
}

func (rs *OurRedSync) NewMutex(key string, options ...redsync.Option) *OurMutex {
	return NewMutex(rs.Redsync.NewMutex(key, options...))
}
