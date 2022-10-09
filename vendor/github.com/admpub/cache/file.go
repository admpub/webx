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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/webx-top/com"

	"github.com/admpub/cache/encoding"
)

var cacheItemPool = sync.Pool{
	New: func() interface{} {
		return &Item{}
	},
}

func CacheItemPoolGet() *Item {
	return cacheItemPool.Get().(*Item)
}

func CacheItemPoolRelease(item *Item) {
	item.Reset()
	cacheItemPool.Put(item)
}

// Item represents a cache item.
type Item struct {
	Val     interface{}
	Created int64
	Expire  int64
}

func (item *Item) hasExpired() bool {
	return item.Expire > 0 &&
		(time.Now().Unix()-item.Created) >= item.Expire
}

func (item *Item) Reset() {
	item.Val = nil
	item.Created = 0
	item.Expire = 0
}

// FileCacher represents a file cache adapter implementation.
type FileCacher struct {
	GetAs
	codec    encoding.Codec
	lock     sync.RWMutex
	rootPath string
	interval int // GC interval.
}

// NewFileCacher creates and returns a new file cacher.
func NewFileCacher() *FileCacher {
	c := &FileCacher{codec: DefaultCodec}
	c.GetAs = GetAs{Cache: c}
	return c
}

func (c *FileCacher) SetCodec(codec encoding.Codec) {
	c.codec = codec
}

func (c *FileCacher) Codec() encoding.Codec {
	return c.codec
}

func (c *FileCacher) filepath(key string) string {
	m := md5.Sum([]byte(key))
	hash := hex.EncodeToString(m[:])
	return filepath.Join(c.rootPath, string(hash[0]), string(hash[1]), hash)
}

// Put puts value into cache with key and expire time.
// If expired is 0, it will be deleted by next GC operation.
func (c *FileCacher) Put(key string, val interface{}, expire int64) error {
	filename := c.filepath(key)

	item := CacheItemPoolGet()
	item.Val = val
	item.Created = time.Now().Unix()
	item.Expire = expire

	data, err := c.codec.Marshal(item)
	CacheItemPoolRelease(item)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filename, data, os.ModePerm)
}

func (c *FileCacher) read(key string, value interface{}) (*Item, error) {
	filename := c.filepath(key)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	item := CacheItemPoolGet()
	item.Val = value

	return item, c.codec.Unmarshal(data, item)
}

// Get gets cached value by given key.
func (c *FileCacher) Get(key string, value interface{}) error {
	item, err := c.read(key, value)
	if item != nil {
		defer CacheItemPoolRelease(item)
	}
	if err != nil {
		return err
	}
	if item.Val == nil {
		return ErrNotFound
	}

	if item.hasExpired() {
		os.Remove(c.filepath(key))
		return ErrExpired
	}
	return nil
}

// Delete deletes cached value by given key.
func (c *FileCacher) Delete(key string) error {
	return os.Remove(c.filepath(key))
}

// Incr increases cached int-type value by given key as a counter.
func (c *FileCacher) Incr(key string) error {
	var i int64
	item, err := c.read(key, &i)
	if item != nil {
		defer CacheItemPoolRelease(item)
	}
	if err != nil {
		return err
	}

	item.Val, err = Incr(i)
	if err != nil {
		return err
	}

	return c.Put(key, item.Val, item.Expire)
}

// Decr cached int value.
func (c *FileCacher) Decr(key string) error {
	var i int64
	item, err := c.read(key, &i)
	if item != nil {
		defer CacheItemPoolRelease(item)
	}
	if err != nil {
		return err
	}

	item.Val, err = Decr(i)
	if err != nil {
		return err
	}

	return c.Put(key, item.Val, item.Expire)
}

// IsExist returns true if cached value exists.
func (c *FileCacher) IsExist(key string) bool {
	return com.IsExist(c.filepath(key))
}

// Flush deletes all cached data.
func (c *FileCacher) Flush() error {
	return os.RemoveAll(c.rootPath)
}

func (c *FileCacher) startGC() {
	c.lock.RLock()
	if c.interval < 1 {
		c.lock.RUnlock()
		return
	}
	c.lock.RUnlock()

	item := &Item{}
	if err := filepath.Walk(c.rootPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk: %v", err)
		}

		if fi.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("ReadFile: %v", err)
		}
		item.Reset()
		if err = c.codec.Unmarshal(data, item); err != nil {
			log.Printf("error garbage collecting cache files: unmarshal: %v", err)
			if err = os.Remove(path); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove: %v", err)
			}
			return nil
		}
		if item.hasExpired() {
			if err = os.Remove(path); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove: %v", err)
			}
		}
		return nil
	}); err != nil {
		log.Printf("error garbage collecting cache files: %v", err)
	}

	time.AfterFunc(time.Duration(c.interval)*time.Second, func() { c.startGC() })
}

// StartAndGC starts GC routine based on config string settings.
func (c *FileCacher) StartAndGC(opt Options) error {
	c.lock.Lock()
	c.rootPath = opt.AdapterConfig
	c.interval = opt.Interval
	var err error
	c.rootPath, err = filepath.Abs(c.rootPath)
	c.lock.Unlock()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(c.rootPath, os.ModePerm); err != nil {
		return err
	}

	go c.startGC()
	return nil
}

func (c *FileCacher) Close() error {
	c.lock.Lock()
	c.interval = 0
	c.lock.Unlock()
	return nil
}

func (c *FileCacher) Client() interface{} {
	return nil
}

func init() {
	Register("file", NewFileCacher())
}
