package agent

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	modelAgent "github.com/admpub/webx/application/model/official/agent"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func levelFormFilter(options ...formfilter.Options) echo.FormDataFilter {
	options = append(
		options,
		formfilter.Exclude(`Created`, `Updated`),
		formfilter.JoinValues(`RoleIds`),
	)
	return formfilter.Build(options...)
}

func LevelIndex(ctx echo.Context) error {
	m := modelAgent.NewAgentLevel(ctx)
	cond := db.Compounds{}
	common.SelectPageCond(ctx, &cond, `id`, `name%`)
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ret := handler.Err(ctx, err)
	list := m.Objects()
	ctx.Set(`listData`, list)
	return ctx.Render(`official/agent/level_index`, ret)
}

func LevelAdd(ctx echo.Context) error {
	var err error
	m := modelAgent.NewAgentLevel(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomerAgentLevel, levelFormFilter())
		if err == nil {
			if len(ctx.FormValues(`roleIds`)) == 0 {
				m.RoleIds = ``
			}
			_, err = m.Add()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/official/agent/level_index`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCustomerAgentLevel, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/official/agent/level_index`)
	ctx.Set(`levelList`, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	ctx.Set(`title`, ctx.T(`添加代理等级`))
	levelSetFormData(ctx, m)
	return ctx.Render(`official/agent/level_edit`, handler.Err(ctx, err))
}

func levelSetFormData(ctx echo.Context, m *modelAgent.AgentLevel) {
	roleM := modelCustomer.NewRole(ctx)
	roleM.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `name`, `description`)
	}, 0, -1, db.And(db.Cond{`parent_id`: 0}))
	ctx.Set(`roleList`, roleM.Objects())

	var roleIds []uint
	if len(m.RoleIds) > 0 {
		roleIds = param.StringSlice(strings.Split(m.RoleIds, `,`)).Uint()
	}
	ctx.SetFunc(`isChecked`, func(roleId uint) bool {
		for _, rid := range roleIds {
			if rid == roleId {
				return true
			}
		}
		return false
	})
}

func LevelEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := modelAgent.NewAgentLevel(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomerAgentLevel, levelFormFilter())
		if err == nil {
			m.Id = id
			if len(ctx.FormValues(`roleIds`)) == 0 {
				m.RoleIds = ``
			}
			err = m.Edit(nil, db.Cond{`id`: id})
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/official/agent/level_index`))
		}
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			if !common.IsBoolFlag(disabled) {
				return ctx.NewError(code.InvalidParameter, ``).SetZone(`disabled`)
			}
			m.Disabled = disabled
			data := ctx.Data()
			err = m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCustomerAgentLevel, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	ctx.Set(`activeURL`, `/official/agent/level_index`)
	ctx.Set(`levelList`, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	ctx.Set(`title`, ctx.T(`修改代理等级`))
	levelSetFormData(ctx, m)
	return ctx.Render(`official/agent/level_edit`, handler.Err(ctx, err))
}

func LevelDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelAgent.NewAgentLevel(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/official/agent/level_index`))
}
