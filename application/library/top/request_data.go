package top

import (
	"github.com/admpub/copier"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type AsMap interface {
	AsMap() param.Store
}

type AsPartialMap interface {
	AsMap(...string) param.Store
}

func RecvValidated(ctx echo.Context, recvs ...factory.Model) error {
	data := GetValidated(ctx)
	switch am := data.(type) {
	case AsMap:
		mp := am.AsMap()
		for _, recv := range recvs {
			recv.Set(mp)
		}
		return nil

	case AsPartialMap:
		mp := am.AsMap()
		for _, recv := range recvs {
			recv.Set(mp)
		}
		return nil

	default:
		for _, recv := range recvs {
			if err := copier.Copy(recv, data); err != nil {
				return err
			}
		}
		return nil
	}
}

func GetValidated(ctx echo.Context) interface{} {
	return ctx.Internal().Get(`validated`)
}
