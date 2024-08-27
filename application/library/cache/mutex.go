package cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var ErrFailedToAcquireLock = errors.New("failed to acquire lock")

var (
	defaultLockType int32 = int32(LockTypeMemory)
	tryLockers            = map[LockType]TryLocker{
		LockTypeMemory: &mutexMemory{},
		LockTypeRedis:  &mutexRedis{},
	}
)

type LockType int32

const (
	LockTypeMemory LockType = iota
	LockTypeRedis
)

func DefaultLockType() LockType {
	v := atomic.LoadInt32(&defaultLockType)
	return LockType(v)
}

func SetDefaultLockType(lockType LockType) {
	atomic.AddInt32(&defaultLockType, int32(lockType))
}

func RegisterTryLocker(t LockType, fn TryLocker) {
	tryLockers[t] = fn
}

type UnlockFunc func() error
type TryLocker interface {
	Lock(key string) (unlock UnlockFunc, err error)
	TryLock(key string) (unlock UnlockFunc, err error)
	TryLockWithTimeout(key string, timeout time.Duration) (unlock UnlockFunc, err error)
	TryLockWithContext(key string, ctx context.Context) (unlock UnlockFunc, err error)
	Forget(key string)
}

func GetLocker(types ...LockType) TryLocker {
	var t LockType
	if len(types) > 0 {
		t = types[0]
	} else {
		t = DefaultLockType()
	}
	if tryLocker, ok := tryLockers[t]; ok {
		return tryLocker
	}
	return tryLockers[LockTypeMemory]
}

func Lock(key string, types ...LockType) (unlock UnlockFunc, err error) {
	return GetLocker(types...).Lock(key)
}

func TryLock(key string, types ...LockType) (unlock UnlockFunc, err error) {
	return GetLocker(types...).TryLock(key)
}

func TryLockWithTimeout(key string, timeout time.Duration, types ...LockType) (unlock UnlockFunc, err error) {
	return GetLocker(types...).TryLockWithTimeout(key, timeout)
}

func TryLockWithContext(key string, ctx context.Context, types ...LockType) (unlock UnlockFunc, err error) {
	return GetLocker(types...).TryLockWithContext(key, ctx)
}

func ForgetLock(key string, types ...LockType) {
	GetLocker(types...).Forget(key)
}
