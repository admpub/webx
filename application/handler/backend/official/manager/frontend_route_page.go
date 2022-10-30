package manager

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/model/official"
)

func FrontendRoutePage(ctx echo.Context) error {
	m := official.NewRoutePage(ctx)
	typ := ctx.Form(`type`)
	cond := db.NewCompounds()
	if len(typ) > 0 {
		cond.AddKV(`page_type`, typ)
	}
	_, err := common.NewLister(m.OfficialCommonRoutePage, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()).Paging(ctx)
	if err != nil {
		return err
	}

	list := m.Objects()
	ctx.Set(`listData`, list)
	ctx.Set(`typeList`, official.RoutePageTypes.Slice())
	ctx.Set(`type`, typ)
	ctx.SetFunc(`typeName`, func(typ string) string {
		return official.RoutePageTypes.Get(typ)
	})
	return ctx.Render(`official/manager/frontend/route_page`, handler.Err(ctx, err))
}

func FrontendRoutePageAdd(ctx echo.Context) error {
	var err error
	m := official.NewRoutePage(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonRoutePage, echo.ExcludeFieldName(`updated`))
		if err == nil {
			m.Method = strings.Join(ctx.FormxValues(`method[]`), `,`)
			_, err = m.Add()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/manager/frontend/route_page`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint64()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCommonRoutePage, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
				if len(m.Method) > 0 {
					for k, v := range strings.Split(m.Method, `,`) {
						if k == 0 {
							ctx.Request().Form().Set(`method[]`, v)
						} else {
							ctx.Request().Form().Add(`method[]`, v)
						}
					}
				}
			}
		}
	}

	ctx.Set(`activeURL`, `/manager/frontend/route_page`)
	ctx.Set(`title`, ctx.T(`添加路由页面`))
	ctx.Set(`typeList`, official.RoutePageTypes.Slice())
	ctx.Set(`methodList`, echo.Methods())
	ctx.SetFunc(`methodChecked`, func(method string, methods []string) bool {
		return com.InSlice(method, methods)
	})
	return ctx.Render(`official/manager/frontend/route_page_edit`, err)
}

func FrontendRoutePageEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint64()
	m := official.NewRoutePage(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonRoutePage, echo.ExcludeFieldName(`created`))
		if err == nil {
			m.Method = strings.Join(ctx.FormxValues(`method[]`), `,`)
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/manager/frontend/route_page`))
			}
		}
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
		echo.StructToForm(ctx, m.OfficialCommonRoutePage, ``, echo.LowerCaseFirstLetter)
		if len(m.Method) > 0 {
			for k, v := range strings.Split(m.Method, `,`) {
				if k == 0 {
					ctx.Request().Form().Set(`method[]`, v)
				} else {
					ctx.Request().Form().Add(`method[]`, v)
				}
			}
		}
	}

	ctx.Set(`activeURL`, `/manager/frontend/route_page`)
	ctx.Set(`title`, ctx.T(`编辑路由页面`))
	ctx.Set(`typeList`, official.RoutePageTypes.Slice())
	ctx.Set(`methodList`, echo.Methods())
	ctx.SetFunc(`methodChecked`, func(method string, methods []string) bool {
		return com.InSlice(method, methods)
	})
	return ctx.Render(`official/manager/frontend/route_page_edit`, err)
}

func FrontendRoutePageDelete(ctx echo.Context) error {
	ids := ctx.FormxValues(`id`).Uint64(func(index int, value uint64) bool {
		return value > 0
	})
	m := official.NewRoutePage(ctx)
	var err error
	for _, _v := range ids {
		if err = m.Delete(nil, db.Cond{`id`: _v}); err != nil {
			break
		}
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/manager/frontend/route_page`))
}
