package cache

import (
	"context"
	"time"

	lock "github.com/admpub/go-lock"
)

var (
	mutexGroup = lock.NewGroup(nil)
)

type mutexMemory struct{}

func (*mutexMemory) Lock(key string) (unlock UnlockFunc, err error) {
	mutexGroup.Lock(key)
	unlock = func() error {
		mutexGroup.UnlockAndFree(key)
		return nil
	}
	return
}

func (*mutexMemory) TryLock(key string) (unlock UnlockFunc, err error) {
	if !mutexGroup.TryLock(key) {
		err = ErrFailedToAcquireLock
		return
	}
	unlock = func() error {
		mutexGroup.UnlockAndFree(key)
		return nil
	}
	return
}

func (*mutexMemory) TryLockWithTimeout(key string, timeout time.Duration) (unlock UnlockFunc, err error) {
	if !mutexGroup.TryLockWithTimeout(key, timeout) {
		err = ErrFailedToAcquireLock
		return
	}
	unlock = func() error {
		mutexGroup.UnlockAndFree(key)
		return nil
	}
	return
}

func (*mutexMemory) TryLockWithContext(key string, ctx context.Context) (unlock UnlockFunc, err error) {
	if !mutexGroup.TryLockWithContext(key, ctx) {
		err = ErrFailedToAcquireLock
		return
	}
	unlock = func() error {
		mutexGroup.UnlockAndFree(key)
		return nil
	}
	return
}

func (*mutexMemory) Forget(key string) {
	mutexGroup.UnlockAndFree(key)
}
