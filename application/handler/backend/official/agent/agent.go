package agent

import (
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"

	modelAgent "github.com/admpub/webx/application/model/official/agent"
)

func formFilter(options ...formfilter.Options) echo.FormDataFilter {
	options = append(
		options,
		formfilter.Exclude(`Created`, `Updated`, `Sold`, `Members`, `RecvMoneyTimes`, `RecvMoneyTotal`, `EarningBalance`, `FreezeAmount`),
	)
	return formfilter.Build(options...)
}

func Index(ctx echo.Context) error {
	m := modelAgent.NewAgentProfile(ctx)
	cond := db.NewCompounds()
	list, err := m.ListPage(cond, `-created`)
	ctx.Set(`listData`, list)
	ctx.Set(`statusList`, modelAgent.ProfileStatus.Slice())
	ctx.Set(`recvMethodList`, modelAgent.RecvMoneyMethod.Slice())
	ctx.SetFunc(`getStatusName`, modelAgent.ProfileStatus.Get)
	return ctx.Render(`official/agent/index`, common.Err(ctx, err))
}

func Add(ctx echo.Context) error {
	var err error
	m := modelAgent.NewAgentProfile(ctx)
	rvM := modelAgent.NewAgentRecv(ctx)
	if ctx.IsPost() {
		ctx.Begin()
		err = ctx.MustBind(m.OfficialCustomerAgentProfile, formFilter())
		if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			err = ctx.MustBind(rvM.OfficialCustomerAgentRecv, formFilter())
			if err == nil {
				_, err = rvM.Add()
			}
		}
		if err == nil {
			ctx.Commit()
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/official/agent/index`))
		}
		ctx.Rollback()
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `customer_id`, id)
			rvM := modelAgent.NewAgentRecv(ctx)
			rvM.Get(nil, db.Cond{`customer_id`: id})
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCustomerAgentProfile, ``, echo.LowerCaseFirstLetter)
				if rvM.CustomerId > 0 {
					echo.StructToForm(ctx, rvM.OfficialCustomerAgentRecv, ``, echo.LowerCaseFirstLetter)
				}
				ctx.Request().Form().Set(`customer_id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/official/agent/index`)
	ctx.Set(`title`, ctx.T(`添加代理商`))
	ctx.Set(`isEdit`, false)
	ctx.Set(`statusList`, modelAgent.ProfileStatus.Slice())
	lvM := modelAgent.NewAgentLevel(ctx)
	lvM.ListByOffset(nil, nil, 0, -1, db.Cond{`disabled`: `N`})
	ctx.Set(`levelList`, lvM.Objects())
	ctx.Set(`recvMethodList`, modelAgent.RecvMoneyMethod.Slice())
	return ctx.Render(`official/agent/edit`, common.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	var err error
	cid := ctx.Formx(`id`).Uint()
	m := modelAgent.NewAgentProfile(ctx)
	err = m.Get(nil, db.Cond{`customer_id`: cid})
	if err != nil {
		return err
	}
	rvM := modelAgent.NewAgentRecv(ctx)
	rvM.Get(nil, db.Cond{`customer_id`: cid})
	if ctx.IsPost() {
		ctx.Begin()
		err = ctx.MustBind(m.OfficialCustomerAgentProfile, formFilter(formfilter.Exclude(`customerId`)))
		if err == nil {
			err = m.Edit(nil, db.Cond{`customer_id`: cid})
		}
		if err == nil {
			err = ctx.MustBind(rvM.OfficialCustomerAgentRecv, formFilter(formfilter.Exclude(`customerId`)))
			if err == nil {
				err = rvM.Edit(nil, db.Cond{`customer_id`: cid})
			}
		}
		if err == nil {
			ctx.Commit()
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/official/agent/level_index`))
		}
		ctx.Rollback()
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCustomerAgentProfile, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
		if rvM.CustomerId > 0 {
			echo.StructToForm(ctx, rvM.OfficialCustomerAgentRecv, ``, echo.LowerCaseFirstLetter)
		}
	}

	ctx.Set(`activeURL`, `/official/agent/index`)
	ctx.Set(`title`, ctx.T(`修改代理商`))
	ctx.Set(`isEdit`, true)
	ctx.Set(`statusList`, modelAgent.ProfileStatus.Slice())
	lvM := modelAgent.NewAgentLevel(ctx)
	lvM.ListByOffset(nil, nil, 0, -1, db.Cond{`disabled`: `N`})
	ctx.Set(`levelList`, lvM.Objects())
	ctx.Set(`recvMethodList`, modelAgent.RecvMoneyMethod.Slice())
	return ctx.Render(`official/agent/edit`, common.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	cid := ctx.Formx(`id`).Uint()
	ctx.Begin()
	m := modelAgent.NewAgentProfile(ctx)
	err := m.Delete(nil, db.Cond{`customer_id`: cid})
	if err == nil {
		rvM := modelAgent.NewAgentRecv(ctx)
		err = rvM.Delete(nil, db.Cond{`customer_id`: cid})
	}
	ctx.End(err == nil)
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/agent/index`))
}
