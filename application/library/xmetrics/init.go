package xmetrics

import "github.com/admpub/nging/v5/application/library/config/extend"

const Name = `metrics`

func init() {
	extend.Register(Name, func() interface{} {
		return New()
	})
}
