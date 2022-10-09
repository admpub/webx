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
	"errors"

	"github.com/webx-top/echo/param"
)

func As(cache Cache) GetAs {
	return GetAs{Cache: cache}
}

type GetAs struct {
	Cache
}

func (g GetAs) String(key string) string {
	var r string
	g.Get(key, &r)
	return r
}

func (g GetAs) Int(key string) int {
	var r int
	g.Get(key, &r)
	return r
}

func (g GetAs) Uint(key string) uint {
	var r uint
	g.Get(key, &r)
	return r
}

func (g GetAs) Int64(key string) int64 {
	var r int64
	g.Get(key, &r)
	return r
}

func (g GetAs) Uint64(key string) uint64 {
	var r uint64
	g.Get(key, &r)
	return r
}

func (g GetAs) Int32(key string) int32 {
	var r int32
	g.Get(key, &r)
	return r
}

func (g GetAs) Uint32(key string) uint32 {
	var r uint32
	g.Get(key, &r)
	return r
}

func (g GetAs) Float32(key string) float32 {
	var r float32
	g.Get(key, &r)
	return r
}

func (g GetAs) Float64(key string) float64 {
	var r float64
	g.Get(key, &r)
	return r
}

func (g GetAs) Bytes(key string) []byte {
	var r []byte
	g.Get(key, &r)
	return r
}

func (g GetAs) Map(key string) map[string]interface{} {
	r := map[string]interface{}{}
	g.Get(key, &r)
	return r
}

func (g GetAs) Mapx(key string) param.Store {
	r := param.Store{}
	g.Get(key, &r)
	return r
}

func (g GetAs) Any(key string) interface{} {
	var r interface{}
	g.Get(key, &r)
	return r
}

func (g GetAs) Slice(key string) []interface{} {
	var r []interface{}
	g.Get(key, &r)
	return r
}

func Incr(val interface{}) (interface{}, error) {
	switch v := val.(type) {
	case int:
		val = v + 1
	case int32:
		val = v + 1
	case int64:
		val = v + 1
	case uint:
		val = v + 1
	case uint32:
		val = v + 1
	case uint64:
		val = v + 1
	default:
		return val, errors.New("item value is not int-type")
	}
	return val, nil
}

func Decr(val interface{}) (interface{}, error) {
	switch v := val.(type) {
	case int:
		val = v - 1
	case int32:
		val = v - 1
	case int64:
		val = v - 1
	case uint:
		if v > 0 {
			val = v - 1
		} else {
			return val, errors.New("item value is less than 0")
		}
	case uint32:
		if v > 0 {
			val = v - 1
		} else {
			return val, errors.New("item value is less than 0")
		}
	case uint64:
		if v > 0 {
			val = v - 1
		} else {
			return val, errors.New("item value is less than 0")
		}
	default:
		return val, errors.New("item value is not int-type")
	}
	return val, nil
}
