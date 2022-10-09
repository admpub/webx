package lock

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Group ...
type Group interface {
	Lock(key interface{})
	TryLock(key interface{}) bool
	TryLockWithTimeout(key interface{}, duration time.Duration) bool
	TryLockWithContext(key interface{}, ctx context.Context) bool
	Unlock(key interface{})
	UnlockAndFree(key interface{})
}

func NewGroup(fn func() Mutex) Group {
	if fn == nil {
		fn = NewChanMutexInterface
	}
	return &group{
		fn:    fn,
		group: make(map[interface{}]*entry),
	}
}

type group struct {
	mu    sync.Mutex
	fn    func() Mutex
	group map[interface{}]*entry
}

type entry struct {
	ref int32
	mu  Mutex
}

func (m *group) get(i interface{}, ref int32, callbacks ...func(*entry)) Mutex {
	m.mu.Lock()
	defer m.mu.Unlock()
	en, ok := m.group[i]
	if !ok {
		if ref > 0 {
			en = &entry{mu: m.fn()}
			m.group[i] = en
		} else {
			return nil
		}
	}
	atomic.AddInt32(&en.ref, ref)
	for _, cb := range callbacks {
		cb(en)
	}
	return en.mu
}

func (m *group) Lock(i interface{}) {
	m.get(i, 1).Lock()
}

func (m *group) Unlock(i interface{}) {
	mu := m.get(i, -1)
	if mu != nil {
		mu.Unlock()
	}
}

func (m *group) UnlockAndFree(i interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	en, ok := m.group[i]
	if !ok {
		return
	}
	ref := atomic.AddInt32(&en.ref, -1)
	if ref < 1 {
		delete(m.group, i)
	}
	en.mu.Unlock()
}

func (m *group) TryLock(i interface{}) (locked bool) {
	m.get(i, 1, func(en *entry) {
		locked = en.mu.TryLock()
		if !locked {
			atomic.AddInt32(&en.ref, -1)
		}
	}).TryLock()
	return
}

func (m *group) TryLockWithTimeout(i interface{}, timeout time.Duration) (locked bool) {
	m.get(i, 1, func(en *entry) {
		locked = en.mu.TryLockWithTimeout(timeout)
		if !locked {
			atomic.AddInt32(&en.ref, -1)
		}
	})
	return
}

func (m *group) TryLockWithContext(i interface{}, ctx context.Context) (locked bool) {
	m.get(i, 1, func(en *entry) {
		locked = en.mu.TryLockWithContext(ctx)
		if !locked {
			atomic.AddInt32(&en.ref, -1)
		}
	})
	return
}
