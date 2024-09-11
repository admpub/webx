package xmetrics

import "github.com/coscms/webcore/library/config/extend"

const Name = `metrics`

func init() {
	extend.Register(Name, func() interface{} {
		return New()
	})
}
