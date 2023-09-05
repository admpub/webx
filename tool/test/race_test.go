package test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type testData struct {
	N int
}

type testDataWithLock struct {
	cfg *testData
	mu  sync.RWMutex
}

func (r *testDataWithLock) SetConfig(cfg *testData) {
	r.mu.Lock()
	r.cfg = cfg
	r.mu.Unlock()
}

func (r *testDataWithLock) Config() *testData {
	r.mu.RLock()
	cfg := r.cfg
	r.mu.RUnlock()
	return cfg
}

func TestRace(t *testing.T) {
	var g = &testDataWithLock{}
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 50; i++ {
		go func() {
			fmt.Printf("%+v\n", g.Config())
			time.Sleep(time.Microsecond * 20)
			wg.Done()
		}()
	}
	for i := 0; i < 50; i++ {
		go func(i int) {
			c := testData{N: i}
			g.SetConfig(&c)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
