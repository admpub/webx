package goredis5

import (
	"context"
	"strings"
	"time"

	redsyncredis "github.com/admpub/redsync/v4/redis"
	"gopkg.in/redis.v5"
)

type Pool struct {
	delegate *redis.Client
}

func (self *Pool) Get(ctx context.Context) (redsyncredis.Conn, error) {
	return &Conn{
		delegate: self.delegate,
		context:  ctx,
	}, nil
}

func NewPool(delegate *redis.Client) *Pool {
	return &Pool{delegate}
}

type Conn struct {
	delegate *redis.Client
	context  context.Context
}

func (self *Conn) Get(name string) (string, error) {
	value, err := self.delegate.Get(name).Result()
	err = noErrNil(err)
	return value, err
}

func (self *Conn) Set(name string, value string) (bool, error) {
	reply, err := self.delegate.Set(name, value, 0).Result()
	return err == nil && reply == "OK", err
}

func (self *Conn) SetNX(name string, value string, expiry time.Duration) (bool, error) {
	return self.delegate.SetNX(name, value, expiry).Result()
}

func (self *Conn) PTTL(name string) (time.Duration, error) {
	return self.delegate.PTTL(name).Result()
}

func (self *Conn) Eval(script *redsyncredis.Script, keysAndArgs ...interface{}) (interface{}, error) {
	var keys []string
	var args []interface{}

	if script.KeyCount > 0 {

		keys = []string{}

		for i := 0; i < script.KeyCount; i++ {
			keys = append(keys, keysAndArgs[i].(string))
		}

		args = keysAndArgs[script.KeyCount:]

	} else {
		keys = []string{}
		args = keysAndArgs
	}

	v, err := self.delegate.EvalSha(script.Hash, keys, args...).Result()
	if err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT ") {
		v, err = self.delegate.Eval(script.Src, keys, args...).Result()
	}
	err = noErrNil(err)
	return v, err
}

func (self *Conn) Close() error {
	// Not needed for this library
	return nil
}

func noErrNil(err error) error {
	if err != nil && err.Error() == "redis: nil" {
		return nil
	}
	return err
}
