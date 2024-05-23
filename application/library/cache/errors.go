package cache

import (
	"errors"

	"github.com/admpub/cache"
)

func IsNotExist(err error) bool {
	return errors.Is(err, cache.ErrNotFound) || errors.Is(err, cache.ErrExpired)
}
