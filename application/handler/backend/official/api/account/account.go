// Package account api 外部接口账号
package account

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official"
	modelApi "github.com/admpub/webx/application/model/official/api"
)

// Index 应用列表
func Index(ctx echo.Context) error {
	var err error
	m := modelApi.NewAccount(ctx)
	cond := db.Compounds{}
	groupId := ctx.Formx(`groupId`).Uint()
	name := ctx.Formx(`q`).String()
	if len(name) > 0 {
		cond.And(db.Or(
			db.Cond{`name`: db.Like(name + `%`)},
			db.Cond{`app_id`: name},
		))
	}
	if groupId > 0 {
		cond.AddKV(`group_id`, groupId)
	}
	sorts := common.Sorts(ctx, `official_common_api_account`, `-id`)
	_, err = common.NewLister(m.OfficialCommonApiAccount, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(ctx)
	ctx.Set(`listData`, m.OfficialCommonApiAccount.Objects())
	mg := official.NewGroup(ctx)
	var groupList []*dbschema.OfficialCommonGroup
	mg.ListByOffset(&groupList, nil, 0, -1, `type`, `api`)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`groupId`, groupId)
	return ctx.Render(`official/api/account/index`, common.Err(ctx, err))
}

func formFilter() echo.FormDataFilter {
	return formfilter.Build(
		formfilter.Exclude(`created`, `updated`),
	)
}

// Add 创建应用
func Add(ctx echo.Context) error {
	user := handler.User(ctx)
	var (
		err error
		id  uint64
	)
	m := modelApi.NewAccount(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(
			m.OfficialCommonApiAccount,
			formFilter(),
		)
		if err != nil {
			goto END
		}
		m.OwnerId = uint64(user.Id)
		m.OwnerType = `user`
		_, err = m.Add()
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`添加成功`))
		return ctx.Redirect(handler.URLFor(`/official/api/account/index`))
	}
	id = ctx.Formx(`copyId`).Uint64()
	if id > 0 {
		err = m.Get(nil, `id`, id)
		if err == nil {
			echo.StructToForm(ctx, m.OfficialCommonApiAccount, ``, echo.LowerCaseFirstLetter)
			ctx.Request().Form().Set(`id`, `0`)
		}
	}

END:
	ctx.Set(`activeURL`, `/official/api/account/index`)
	ctx.Set(`title`, ctx.T(`添加账号`))
	mg := official.NewGroup(ctx)
	var groupList []*dbschema.OfficialCommonGroup
	mg.ListByOffset(&groupList, nil, 0, -1, `type`, `api`)
	ctx.Set(`groupList`, groupList)
	return ctx.Render(`official/api/account/edit`, common.Err(ctx, err))
}

// Edit 修改应用
func Edit(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.NewError(code.InvalidParameter, `参数“%s”值无效`, `id`).SetZone(`id`)
	}
	m := modelApi.NewAccount(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `账号不存在`)
		}
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(
			m.OfficialCommonApiAccount,
			formFilter(),
		)
		if err != nil {
			goto END
		}
		err = m.Edit(nil, `id`, m.Id)
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`修改成功`))
		return ctx.Redirect(handler.URLFor(`/official/api/account/index`))
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
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
		echo.StructToForm(ctx, m.OfficialCommonApiAccount, ``, echo.LowerCaseFirstLetter)
	}

END:
	ctx.Set(`activeURL`, `/official/api/account/index`)
	ctx.Set(`title`, ctx.T(`修改账号`))
	mg := official.NewGroup(ctx)
	var groupList []*dbschema.OfficialCommonGroup
	mg.ListByOffset(&groupList, nil, 0, -1, `type`, `api`)
	ctx.Set(`groupList`, groupList)
	return ctx.Render(`official/api/account/edit`, common.Err(ctx, err))
}

// Delete 删除应用
func Delete(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.NewError(code.InvalidParameter, `参数“%s”值无效`, `id`).SetZone(`id`)
	}
	m := modelApi.NewAccount(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `账号不存在`)
		}
		return err
	}
	err = m.Delete(nil, `id`, id)
	if err != nil {
		return err
	}
	common.SendOk(ctx, ctx.T(`删除成功`))
	return ctx.Redirect(handler.URLFor(`/official/api/account/index`))
}
