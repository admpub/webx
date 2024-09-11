package download

import (
	"time"

	imageproxy "github.com/admpub/imageproxy"
	"github.com/webx-top/echo"
	"github.com/webx-top/image"

	"github.com/coscms/webcore/library/notice"
	uploadPrepare "github.com/coscms/webcore/registry/upload/prepare"
)

type Config struct {
	echo.Context
	ID            string
	FileURL       string
	Checkin       bool
	Watermark     *image.WatermarkOptions
	Crop          *imageproxy.Options
	PrepareData   *uploadPrepare.PrepareData
	NoticeSender  notice.Noticer
	Progress      *notice.Progress
	MaxMB         int64
	MaxRetries    int
	RetryInterval time.Duration
	DisableChunk  bool
}

var (
	MaxRetries    = 3
	RetryInterval = time.Second
)

func NewConfig(ctx echo.Context) *Config {
	return &Config{
		Context:       ctx,
		MaxMB:         0,
		MaxRetries:    MaxRetries,
		RetryInterval: RetryInterval,
	}
}

type Options func(*Config)

func OptionsFileURL(fileURL string) Options {
	return func(c *Config) {
		c.FileURL = fileURL
	}
}

func OptionsID(id string) Options {
	return func(c *Config) {
		c.ID = id
	}
}

func OptionsMaxMB(maxMB int64) Options {
	return func(c *Config) {
		c.MaxMB = maxMB
	}
}

func OptionsCheckin(on bool) Options {
	return func(c *Config) {
		c.Checkin = on
	}
}

func OptionsWatermark(opt *image.WatermarkOptions) Options {
	return func(c *Config) {
		c.Watermark = opt
	}
}

func OptionsCrop(opt *imageproxy.Options) Options {
	return func(c *Config) {
		c.Crop = opt
	}
}

func OptionsPrepareData(data *uploadPrepare.PrepareData) Options {
	return func(c *Config) {
		c.PrepareData = data
	}
}

func OptionsNoticeSender(noticeSender notice.Noticer) Options {
	return func(c *Config) {
		c.NoticeSender = noticeSender
	}
}

func OptionsProgress(pro *notice.Progress) Options {
	return func(c *Config) {
		c.Progress = pro
	}
}

func OptionsDisableChunk(disableChunk bool) Options {
	return func(c *Config) {
		c.DisableChunk = disableChunk
	}
}
