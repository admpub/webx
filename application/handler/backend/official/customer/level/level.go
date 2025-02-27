package level

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/nsql"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	modelLevel "github.com/coscms/webfront/model/official/level"
)

func Index(ctx echo.Context) error {
	group := ctx.Form(`group`)
	m := modelLevel.NewLevel(ctx)
	cond := db.NewCompounds()
	if len(group) > 0 {
		cond.AddKV(`group`, group)
	}
	nsql.SelectPageCond(ctx, cond, `id`, `name%`)
	_, err := common.PagingWithLister(ctx, common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ret := common.Err(ctx, err)
	list := m.Objects()
	ctx.Set(`listData`, list)
	ctx.Set(`groupList`, modelLevel.GroupList.Slice())
	ctx.SetFunc(`levelGroupName`, modelLevel.GroupList.Get)
	ctx.SetFunc(`assetTypeName`, modelCustomer.AssetTypes.Get)
	ctx.SetFunc(`amountTypeName`, modelLevel.AmountTypes.Get)
	return ctx.Render(`official/customer/level/index`, ret)
}

func formFilter(options ...formfilter.Options) echo.FormDataFilter {
	options = append(
		options,
		formfilter.Exclude(`Created`, `Updated`),
		formfilter.JoinValues(`RoleIds`),
	)
	return formfilter.Build(options...)
}

func Add(ctx echo.Context) error {
	var err error
	m := modelLevel.NewLevel(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomerLevel, formFilter())
		if err == nil {
			if len(ctx.FormValues(`roleIds`)) == 0 {
				m.RoleIds = ``
			}
			_, err = m.Add()
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/customer/level/index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCustomerLevel, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/official/customer/level/index`)
	ctx.Set(`title`, ctx.T(`添加等级`))
	setFormData(ctx, m)
	return ctx.Render(`official/customer/level/edit`, common.Err(ctx, err))
}

func setFormData(ctx echo.Context, m *modelLevel.Level) {
	ctx.Set(`groupList`, modelLevel.GroupList.Slice())
	ctx.Set(`assetTypes`, modelCustomer.AssetTypes.Slice())
	ctx.Set(`amountTypes`, modelLevel.AmountTypes.Slice())

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

func Edit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := modelLevel.NewLevel(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomerLevel, formFilter())
		if err == nil {
			m.Id = id
			if len(ctx.FormValues(`roleIds`)) == 0 {
				m.RoleIds = ``
			}
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/customer/level/index`))
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

	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCustomerLevel, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	ctx.Set(`activeURL`, `/official/customer/level/index`)
	ctx.Set(`title`, ctx.T(`编辑等级`))
	setFormData(ctx, m)
	return ctx.Render(`official/customer/level/edit`, common.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelLevel.NewLevel(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/customer/level/index`))
}
