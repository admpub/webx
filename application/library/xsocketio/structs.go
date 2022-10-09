package xsocketio

import "github.com/webx-top/echo"

func NewConfig() *Config {
	return &Config{}
}

type Config struct {
}

func (c *Config) FromStore(v echo.H) *Config {
	return c
}
