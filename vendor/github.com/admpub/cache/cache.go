// Copyright 2018 The go-cache Authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cache

import (
	"context"
	"fmt"

	"github.com/webx-top/echo/param"

	"github.com/admpub/cache/encoding"
	"github.com/admpub/cache/encoding/json"
	"github.com/admpub/ini"
)

const _VERSION = "0.2.2"

func Version() string {
	return _VERSION
}

// Cache is the interface that operates the cache data.
type Cache interface {
	// Put puts value into cache with key and expire time.
	Put(ctx context.Context, key string, val interface{}, timeout int64) error
	// Get gets cached value by given key.
	Get(ctx context.Context, key string, value interface{}) error
	// Delete deletes cached value by given key.
	Delete(ctx context.Context, key string) error
	// Incr increases cached int-type value by given key as a counter.
	Incr(ctx context.Context, key string) error
	// Decr decreases cached int-type value by given key as a counter.
	Decr(ctx context.Context, key string) error
	// IsExist returns true if cached value exists.
	IsExist(ctx context.Context, key string) (bool, error)
	// Flush deletes all cached data.
	Flush(ctx context.Context) error
	// StartAndGC starts GC routine based on config string settings.
	StartAndGC(ctx context.Context, opt Options) error
	Close() error
	Client() interface{}
	Codec
	Getter
}

type Getter interface {
	String(ctx context.Context, key string) string
	Int(ctx context.Context, key string) int
	Uint(ctx context.Context, key string) uint
	Int64(ctx context.Context, key string) int64
	Uint64(ctx context.Context, key string) uint64
	Int32(ctx context.Context, key string) int32
	Uint32(ctx context.Context, key string) uint32
	Float32(ctx context.Context, key string) float32
	Float64(ctx context.Context, key string) float64
	Bytes(ctx context.Context, key string) []byte
	Map(ctx context.Context, key string) map[string]interface{}
	Any(ctx context.Context, key string) interface{}
	Mapx(ctx context.Context, key string) param.Store
	Slice(ctx context.Context, key string) []interface{}
}

var DefaultCodec encoding.Codec = json.JSON

type Codec interface {
	SetCodec(encoding.Codec)
	Codec() encoding.Codec
}

// Options represents a struct for specifying configuration options for the cache middleware.
type Options struct {
	// Name of adapter. Default is "memory".
	Adapter string
	// Adapter configuration, it's corresponding to adapter.
	AdapterConfig string
	// GC interval time in seconds. Default is 60.
	Interval int
	// Occupy entire database. Default is false.
	OccupyMode bool
	// Configuration section name. Default is "cache".
	Section string
}

var cfg *ini.File

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}
	if len(opt.Section) == 0 {
		opt.Section = "cache"
	}
	sec := Config().Section(opt.Section)

	if len(opt.Adapter) == 0 {
		opt.Adapter = sec.Key("ADAPTER").MustString("memory")
	}
	if opt.Interval == 0 {
		opt.Interval = sec.Key("INTERVAL").MustInt(60)
	}
	if len(opt.AdapterConfig) == 0 {
		opt.AdapterConfig = sec.Key("ADAPTER_CONFIG").MustString("data/caches")
	}

	return opt
}

// NewCacher creates and returns a new cacher by given adapter name and configuration.
// It panics when given adapter isn't registered and starts GC automatically.
func NewCacher(ctx context.Context, name string, opt Options) (Cache, error) {
	adapter, ok := adapters[name]
	if !ok {
		return nil, fmt.Errorf("cache: unknown adapter '%s'(forgot to import?)", name)
	}
	return adapter, adapter.StartAndGC(ctx, opt)
}

// Cacher is a middleware that maps a cache.Cache service into the Macaron handler chain.
// An single variadic cache.Options struct can be optionally provided to configure.
func Cacher(ctx context.Context, options ...Options) (Cache, error) {
	opt := prepareOptions(options)
	return NewCacher(ctx, opt.Adapter, opt)
}

var adapters = make(map[string]Cache)

// Register registers a adapter.
func Register(name string, adapter Cache) {
	if adapter == nil {
		panic("cache: cannot register adapter with nil value")
	}
	if _, dup := adapters[name]; dup {
		panic(fmt.Errorf("cache: cannot register adapter '%s' twice", name))
	}
	adapters[name] = adapter
}

func Adapters() []string {
	var r []string
	for name := range adapters {
		r = append(r, name)
	}
	return r
}

func HasAdapter(name string) bool {
	_, ok := adapters[name]
	return ok
}

func Config() *ini.File {
	if cfg == nil {
		return ini.Empty()
	}
	return cfg
}
