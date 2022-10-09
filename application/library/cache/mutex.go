package cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var ErrFailedToAcquireLock = errors.New("failed to acquire lock")

var (
	defaultLockType = LockTypeMemory
	tryLockers      = map[int32]TryLocker{
		LockTypeMemory: &mutexMemory{},
		LockTypeRedis:  &mutexRedis{},
	}
)

const (
	LockTypeMemory int32 = iota
	LockTypeRedis
)

func DefaultLockType() int32 {
	return atomic.LoadInt32(&defaultLockType)
}

func SetDefaultLockType(lockType int32) {
	atomic.AddInt32(&defaultLockType, lockType)
}

func RegisterTryLocker(t int32, fn TryLocker) {
	tryLockers[t] = fn
}

type UnlockFunc func() error
type TryLocker interface {
	Lock(key string) (unlock UnlockFunc, err error)
	TryLock(key string) (unlock UnlockFunc, err error)
	TryLockWithTimeout(key string, timeout time.Duration) (unlock UnlockFunc, err error)
	TryLockWithContext(key string, ctx context.Context) (unlock UnlockFunc, err error)
}

func Lock(key string, types ...int32) (unlock UnlockFunc, err error) {
	var t int32
	if len(types) > 0 {
		t = types[0]
	} else {
		t = DefaultLockType()
	}
	if tryLocker, ok := tryLockers[t]; ok {
		return tryLocker.Lock(key)
	}
	return tryLockers[LockTypeMemory].Lock(key)
}

func TryLock(key string, types ...int32) (unlock UnlockFunc, err error) {
	var t int32
	if len(types) > 0 {
		t = types[0]
	} else {
		t = DefaultLockType()
	}
	if tryLocker, ok := tryLockers[t]; ok {
		return tryLocker.TryLock(key)
	}
	return tryLockers[LockTypeMemory].TryLock(key)
}

func TryLockWithTimeout(key string, timeout time.Duration, types ...int32) (unlock UnlockFunc, err error) {
	var t int32
	if len(types) > 0 {
		t = types[0]
	} else {
		t = DefaultLockType()
	}
	if tryLocker, ok := tryLockers[t]; ok {
		return tryLocker.TryLockWithTimeout(key, timeout)
	}
	return tryLockers[LockTypeMemory].TryLockWithTimeout(key, timeout)
}

func TryLockWithContext(key string, ctx context.Context, types ...int32) (unlock UnlockFunc, err error) {
	var t int32
	if len(types) > 0 {
		t = types[0]
	} else {
		t = DefaultLockType()
	}
	if tryLocker, ok := tryLockers[t]; ok {
		return tryLocker.TryLockWithContext(key, ctx)
	}
	return tryLockers[LockTypeMemory].TryLockWithContext(key, ctx)
}
