package xsocketio

import (
	"sync"

	esi "github.com/webx-top/echo-socket.io"
)

var (
	instances = map[string]*esi.Wrapper{}
	mu        = &sync.RWMutex{}
)

func SocketIO(namespace string) *esi.Wrapper {
	mu.RLock()
	v, y := instances[namespace]
	mu.RUnlock()
	if y {
		return v
	}

	v = socketIOWrapper(namespace)
	v.Serve()
	mu.Lock()
	instances[namespace] = v
	mu.Unlock()
	return v
}

func Close(namespace string) bool {
	mu.RLock()
	v, y := instances[namespace]
	mu.RUnlock()
	if y {
		v.Close()
		mu.Lock()
		delete(instances, namespace)
		mu.Unlock()
	}

	return y
}

func CloseAll() {
	mu.Lock()
	for namespace, instance := range instances {
		instance.Close()
		delete(instances, namespace)
	}
	mu.Unlock()
}
