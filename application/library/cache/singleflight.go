package cache

import (
	"golang.org/x/sync/singleflight"
)

type SinglefightResult = singleflight.Result

var _ Singleflighter = (*singleflight.Group)(nil)

type Singleflighter interface {
	Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)
	DoChan(key string, fn func() (interface{}, error)) <-chan SinglefightResult
	Forget(key string)
}

type singleflightLock struct {
	mutex TryLocker
}

func (s *singleflightLock) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	var unlock UnlockFunc
	unlock, err = s.mutex.TryLock(key)
	if err != nil {
		if err == ErrFailedToAcquireLock {
			err = nil
			shared = true
		}
		return
	}
	defer unlock()
	v, err = fn()
	return
}

func (s *singleflightLock) DoChan(key string, fn func() (interface{}, error)) <-chan SinglefightResult {
	ch := make(chan SinglefightResult, 1)
	var unlock UnlockFunc
	unlock, err := s.mutex.TryLock(key)
	if err != nil {
		r := SinglefightResult{}
		if err == ErrFailedToAcquireLock {
			r.Shared = true
		} else {
			r.Err = err
		}
		ch <- r
		return ch
	}
	go func() {
		defer unlock()
		r := SinglefightResult{}
		r.Val, r.Err = fn()
		ch <- r
	}()
	return ch
}

func (s *singleflightLock) Forget(key string) {
	s.mutex.Forget(key)
}

func NewSingleflight(mu TryLocker) Singleflighter {
	return &singleflightLock{mutex: mu}
}

func Singleflight(types ...LockType) Singleflighter {
	return NewSingleflight(GetLocker(types...))
}
