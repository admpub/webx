package official

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

type ClickFlowTargetFunc func(ctx echo.Context, id interface{}) (after func(typ string) error, idGetter func() uint64, err error)

func (c ClickFlowTargetFunc) Do(ctx echo.Context, id interface{}) (after func(typ string) error, idGetter func() uint64, err error) {
	return c(ctx, id)
}

type ClickFlowTarget interface {
	Do(ctx echo.Context, id interface{}) (after func(typ string) error, idGetter func() uint64, err error)
}

var ClickFlowTargets = map[string]ClickFlowTarget{}

func AddClickFlowTarget(name string, target ClickFlowTarget) {
	ClickFlowTargets[name] = target
}

func NewClickFlow(ctx echo.Context) *ClickFlow {
	return &ClickFlow{
		OfficialCommonClickFlow: dbschema.NewOfficialCommonClickFlow(ctx),
	}
}

type ClickFlow struct {
	*dbschema.OfficialCommonClickFlow
}

func (f *ClickFlow) Add() (pk interface{}, err error) {
	var exists bool
	exists, err = f.Exists(f.TargetType, f.TargetId, f.OwnerId, f.OwnerType)
	if err != nil {
		return
	}
	if exists {
		err = f.Context().E(`您已经表过态了`)
		return
	}
	return f.OfficialCommonClickFlow.Insert()
}

func (f *ClickFlow) Exists(targetType string, targetID uint64, ownerID uint64, ownerType string) (bool, error) {
	return f.OfficialCommonClickFlow.Exists(nil, db.And(
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
		db.Cond{`owner_id`: ownerID},
		db.Cond{`owner_type`: ownerType},
	))
}

func (f *ClickFlow) Find(targetType string, targetID uint64, ownerID uint64, ownerType string) error {
	return f.Get(nil, db.And(
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
		db.Cond{`owner_id`: ownerID},
		db.Cond{`owner_type`: ownerType},
	))
}

func (f *ClickFlow) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCommonClickFlow.Update(mw, args...)
}

func (f *ClickFlow) DelByTarget(targetType string, targetID uint64) error {
	return f.OfficialCommonClickFlow.Delete(nil, db.And(
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
	))
}
