# go-cache
This project encapsulates multiple db servers, redis、ledis、memcache、file、memory、nosql、postgresql

example
```go
package main

import (
	"context"

	"github.com/admpub/cache"
	_ "github.com/admpub/cache/redis"
)

func main() {
	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ca, err := cache.Cacher(rootCtx, cache.Options{
		Adapter:       "redis",
		AdapterConfig: "addr=127.0.0.1:6379",
		OccupyMode:    true,
	})

	if err != nil {
		panic(err)
	}

	reqCtx := context.Background()
	ca.Put(reqCtx, "key", "cache", 60)
}
```
