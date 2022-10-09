package cache

import (
	"github.com/admpub/cache/x"
	"github.com/webx-top/db/lib/factory"
)

var DBCacher = NewDBCacher()

func NewDBCacher() factory.Cacher {
	return &dbCacher{}
}

type dbCacher struct{}

func (d *dbCacher) Put(key string, value interface{}, ttlSeconds int64) error {
	return Cache().Put(key, value, ttlSeconds)
}

func (d *dbCacher) Del(key string) error {
	return Cache().Delete(key)
}

func (d *dbCacher) Get(key string, value interface{}) error {
	return Cache().Get(key, value)
}

func (d *dbCacher) Do(key string, recv interface{}, fn func() error, ttlSeconds int64) error {
	return XQuery(key, recv, x.QueryFunc(fn), TTL(ttlSeconds))
}
