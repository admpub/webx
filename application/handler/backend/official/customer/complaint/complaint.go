package complaint

import (
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"

	modelCustomer "github.com/coscms/webfront/model/official/customer"
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
	return ctx.Render(`official/customer/complaint/index`, common.Err(ctx, err))
}

func formFilter() echo.FormDataFilter {
	return formfilter.Build(
		formfilter.Exclude(`updated`, `created`, `content`, `targetIdent`, `targetId`, `targetType`, `type`),
	)
}

func Edit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint64()
	m := modelCustomer.NewComplaint(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonComplaint, formFilter())
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/customer/complaint/index`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonComplaint, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/official/customer/complaint/index`)
	ctx.Set(`types`, modelCustomer.ComplaintTypeList())
	ctx.Set(`targets`, modelCustomer.ComplaintTargetList())
	ctx.Set(`data`, m.OfficialCommonComplaint)
	return ctx.Render(`official/customer/complaint/edit`, err)
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewComplaint(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/customer/complaint/index`))
}
