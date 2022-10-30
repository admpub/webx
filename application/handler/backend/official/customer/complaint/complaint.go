package complaint

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func Index(ctx echo.Context) error {
	m := modelCustomer.NewComplaint(ctx)
	cond := db.NewCompounds()
	t := ctx.Form(`type`)
	if len(t) > 0 {
		cond.AddKV(`type`, t)
	}
	target := ctx.Form(`target`)
	if len(target) > 0 {
		cond.AddKV(`target_type`, target)
	}
	process := ctx.Form(`process`)
	if len(process) > 0 {
		cond.AddKV(`process`, process)
	}
	list, err := m.ListPage(cond, `-id`)
	ctx.Set(`listData`, list)
	ctx.Set(`types`, modelCustomer.ComplaintTypeList())
	ctx.Set(`targets`, modelCustomer.ComplaintTargetList())
	ctx.Set(`processes`, modelCustomer.ComplaintProcessList())
	ctx.Set(`type`, t)
	ctx.Set(`target`, target)
	ctx.Set(`process`, process)
	ctx.SetFunc(`processName`, modelCustomer.ComplaintProcesses.Get)
	return ctx.Render(`official/customer/complaint/index`, handler.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint64()
	m := modelCustomer.NewComplaint(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonComplaint, echo.ExcludeFieldName(`created`))
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/official/customer/complaint/index`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonComplaint, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/official/customer/complaint/index`)
	ctx.Set(`types`, modelCustomer.ComplaintTypeList())
	ctx.Set(`targets`, modelCustomer.ComplaintTargetList())
	return ctx.Render(`official/customer/complaint/edit`, err)
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewComplaint(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/official/customer/complaint/index`))
}
