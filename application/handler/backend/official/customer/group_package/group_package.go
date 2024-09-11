package group_package

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"

	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	modelLevel "github.com/admpub/webx/application/model/official/level"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
)

func Index(ctx echo.Context) error {
	group := ctx.Form(`group`)
	m := modelCustomer.NewGroupPackage(ctx)
	cond := db.NewCompounds()
	if len(group) > 0 {
		cond.AddKV(`group`, group)
	}
	common.SelectPageCond(ctx, cond, `id`, `title`)
	_, err := common.PagingWithLister(ctx, common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ret := common.Err(ctx, err)
	list := m.Objects()
	ctx.Set(`listData`, list)
	ctx.Set(`groupList`, modelLevel.GroupList.Slice())
	ctx.SetFunc(`levelGroupName`, modelLevel.GroupList.Get)
	ctx.SetFunc(`timeUnitName`, modelCustomer.GroupPackageTimeUnits.Get)
	return ctx.Render(`official/customer/group_package/index`, ret)
}

func formFilter(options ...formfilter.Options) echo.FormDataFilter {
	options = append(
		options,
		formfilter.Exclude(`Created`, `Updated`, `Sold`),
	)
	return formfilter.Build(options...)
}

func Add(ctx echo.Context) error {
	var err error
	m := modelCustomer.NewGroupPackage(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomerGroupPackage, formFilter())
		if err == nil {
			_, err = m.Add()
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/customer/group_package/index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCustomerGroupPackage, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/official/customer/group_package/index`)
	ctx.Set(`title`, ctx.T(`添加套餐`))
	setFormData(ctx, m)
	return ctx.Render(`official/customer/group_package/edit`, common.Err(ctx, err))
}

func setFormData(ctx echo.Context, m *modelCustomer.GroupPackage) {
	ctx.Set(`groupList`, modelLevel.GroupList.Slice())
	ctx.Set(`timeUnits`, modelCustomer.GroupPackageTimeUnits.Slice())
}

func Edit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewGroupPackage(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomerGroupPackage, formFilter())
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/customer/group_package/index`))
			}
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

		recommend := ctx.Query(`recommend`)
		if len(recommend) > 0 {
			if !common.IsBoolFlag(recommend) {
				return ctx.NewError(code.InvalidParameter, ``).SetZone(`recommend`)
			}
			m.Recommend = recommend
			data := ctx.Data()
			err = m.UpdateField(nil, `recommend`, recommend, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}

	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCustomerGroupPackage, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	ctx.Set(`activeURL`, `/official/customer/group_package/index`)
	ctx.Set(`title`, ctx.T(`编辑套餐`))
	setFormData(ctx, m)
	return ctx.Render(`official/customer/group_package/edit`, common.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewGroupPackage(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/customer/group_package/index`))
}
