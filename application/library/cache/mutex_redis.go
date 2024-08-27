package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/admpub/once"
	"github.com/admpub/redsync/v4"
	goredis "github.com/admpub/redsync/v4/redis/goredis/v5"
	"gopkg.in/redis.v5"
)

var (
	redsyncClient   *redsync.Redsync
	redsyncOnce     once.Once
	maxLockDuration = 2 * time.Minute
)

func resetRedsync() {
	redsyncOnce.Reset()
}

func initRedsync() {
	client, _ := Cache(cacheRootContext, `locker`).Client().(*redis.Client)
	if client == nil {
		client = RedisClient()
	}
	pool := goredis.NewPool(client)
	redsyncClient = redsync.New(pool)
}

func onceInitRedsync() {
	initRedsync()
}

func RedsyncClient() *redsync.Redsync {
	redsyncOnce.Do(onceInitRedsync)

	return redsyncClient
}

// RedisMutex 分布式锁
// example:
// mutex := RedisMutex(`goods_1`)
// err = mutex.Lock(ctx)
//
//	if err != nil {
//		panic(err)
//	}
//
// mutex.Unlock(ctx)
func RedisMutex(key string, options ...redsync.Option) *redsync.Mutex {
	return RedsyncClient().NewMutex(key, options...)
}

type mutexRedis struct{}

func (*mutexRedis) Lock(key string) (unlock UnlockFunc, err error) {
	delay := 100 * time.Millisecond
	m := RedisMutex(key,
		redsync.WithExpiry(maxLockDuration),
		redsync.WithTries(1000),
		redsync.WithRetryDelayFunc(func(tries int) time.Duration {
			return delay * time.Duration(tries)
		}),
		redsync.WithRetryDelay(delay),
	)

	err = m.Lock()
	if err != nil {
		if err == redsync.ErrFailed {
			err = ErrFailedToAcquireLock
		}
		return
	}
	unlock = func() error {
		ok, err := m.Unlock()
		if !ok || err != nil {
			return fmt.Errorf("unlock unsuccessful: %w", err)
		}
		return nil
	}
	return
}

func (*mutexRedis) TryLock(key string) (unlock UnlockFunc, err error) {
	m := RedisMutex(key,
		redsync.WithExpiry(maxLockDuration),
		redsync.WithTries(1),
		redsync.WithRetryDelay(50*time.Millisecond),
	)
	err = m.Lock()
	if err != nil {
		if err == redsync.ErrFailed {
			err = ErrFailedToAcquireLock
		}
		return
	}
	ticker := time.NewTicker(maxLockDuration / 3)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				m.Extend()
			}
		}
	}()

	unlock = func() error {
		close(done)
		ticker.Stop()
		ok, err := m.Unlock()
		if !ok || err != nil {
			return fmt.Errorf("unlock unsuccessful: %w", err)
		}
		return nil
	}
	return
}

func (*mutexRedis) TryLockWithTimeout(key string, maxLockDuration time.Duration) (unlock UnlockFunc, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxLockDuration)
	defer cancel()
	m := RedisMutex(key,
		redsync.WithExpiry(maxLockDuration),
		redsync.WithTries(1),
		redsync.WithRetryDelay(50*time.Millisecond),
	)
	err = m.LockContext(ctx)
	if err != nil {
		if err == redsync.ErrFailed {
			err = ErrFailedToAcquireLock
		}
		return
	}
	ticker := time.NewTicker(maxLockDuration / 3)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				m.Extend()
			}
		}
	}()

	unlock = func() error {
		close(done)
		ticker.Stop()
		ok, err := m.Unlock()
		if !ok || err != nil {
			return fmt.Errorf("unlock unsuccessful: %w", err)
		}
		return nil
	}
	return
}

func (*mutexRedis) TryLockWithContext(key string, ctx context.Context) (unlock UnlockFunc, err error) {
	m := RedisMutex(key,
		redsync.WithExpiry(maxLockDuration),
		redsync.WithTries(1),
		redsync.WithRetryDelay(50*time.Millisecond),
	)
	err = m.LockContext(ctx)
	if err != nil {
		if err == redsync.ErrFailed {
			err = ErrFailedToAcquireLock
		}
		return
	}
	ticker := time.NewTicker(maxLockDuration / 3)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				m.Extend()
			}
		}
	}()

	unlock = func() error {
		close(done)
		ticker.Stop()
		ok, err := m.Unlock()
		if !ok || err != nil {
			return fmt.Errorf("unlock unsuccessful: %w", err)
		}
		return nil
	}
	return
}

func (*mutexRedis) Forget(key string) {
	RedisMutex(key).Unlock()
}
