package customer

import (
	"context"

	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

var (
	ComplaintTypes   = echo.NewKVData()
	ComplaintTargets = echo.NewKVData()
)

type ComplaintFunc func(*dbschema.OfficialCommonComplaint) error

func init() {
	ComplaintTypeAdd(`illegal`, `非法内容`, nil)
	ComplaintTypeAdd(`infringing`, `侵权内容`, nil)
	ComplaintTypeAdd(`other`, `其它`, nil)
}

func ComplaintTypeAdd(key, name string, fn ComplaintFunc) {
	item := &echo.KV{
		K: key,
		V: name,
	}
	if fn != nil {
		item.SetFn(func(context.Context) interface{} {
			return fn
		})
	}
	ComplaintTypes.AddItem(item)
}

func ExecComplaintTypeFunc(key string, mdl *dbschema.OfficialCommonComplaint) error {
	item := ComplaintTypes.GetItem(key)
	if item == nil || item.Fn() == nil {
		return nil
	}
	return item.Fn()(mdl.Context()).(ComplaintFunc)(mdl)
}

func ComplaintTypeList() []*echo.KV {
	return ComplaintTypes.Slice()
}

func ComplaintTargetAdd(key, name string, fn ComplaintFunc, urlFormat func(*dbschema.OfficialCommonComplaint) string) {
	item := &echo.KV{
		K: key,
		V: name,
	}
	if fn != nil {
		item.SetFn(func(context.Context) interface{} {
			return fn
		})
	}
	extraX := echo.H{
		`urlFormat`: urlFormat,
	}
	item.X = extraX
	ComplaintTargets.AddItem(item)
}

func ExecComplaintTargetFunc(key string, mdl *dbschema.OfficialCommonComplaint) error {
	item := ComplaintTargets.GetItem(key)
	if item == nil || item.Fn() == nil {
		return nil
	}
	return item.Fn()(mdl.Context()).(ComplaintFunc)(mdl)
}

func ComplaintTargetList() []*echo.KV {
	return ComplaintTargets.Slice()
}
