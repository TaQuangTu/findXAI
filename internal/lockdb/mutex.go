package lockdb

import (
	"context"

	"github.com/go-redsync/redsync/v4"
)

type OurMutex struct {
	*redsync.Mutex
	ctx context.Context
}

func NewMutex(mux *redsync.Mutex) *OurMutex {
	return &OurMutex{
		Mutex: mux,
	}
}

func (m *OurMutex) TryLockContext(ctx context.Context) (err error) {
	err = m.Mutex.TryLockContext(ctx)
	if err != nil {
		return
	}
	m.ctx = ctx
	return nil
}

func (m *OurMutex) LockContext(ctx context.Context) (err error) {
	err = m.Mutex.LockContext(ctx)
	if err != nil {
		return
	}
	m.ctx = ctx
	return nil
}

func (m *OurMutex) UnlockContext(ctx context.Context) (ok bool, err error) {
	ok, err = m.Mutex.UnlockContext(ctx)
	if err != nil {
		return
	}
	return ok, nil
}

func (m *OurMutex) Unlock() (bool, error) {
	return m.UnlockContext(m.ctx)
}
