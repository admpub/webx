package cache

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrExpired      = errors.New("expired")
	ErrNotSupported = errors.New("not supported operation")
)
