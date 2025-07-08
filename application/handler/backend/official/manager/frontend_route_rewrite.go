package manager

import (
	"strings"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/coscms/webfront/model/official"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

func FrontendRouteRewrite(ctx echo.Context) error {
	m := official.NewRouteRewrite(ctx)
	cond := db.NewCompounds()
	_, err := common.NewLister(m.OfficialCommonRouteRewrite, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()).Paging(ctx)
	if err != nil {
		return err
	}

	list := m.Objects()
	ctx.Set(`listData`, list)
	return ctx.Render(`official/manager/frontend/route_rewrite`, common.Err(ctx, err))
}

func FrontendRouteList(ctx echo.Context) error {
	prefix := ctx.Form(`prefix`)
	size := ctx.Formx(`size`).Int()
	if size < 1 {
		size = 10
	}
	var result []string
	if len(prefix) > 0 {
		routes := frontend.IRegister().Routes()
		var i int
		for _, route := range routes {
			if strings.HasPrefix(route.Path, prefix) {
				result = append(result, route.Path)
				i++
				if i >= size {
					break
				}
			}
		}
	}
	data := ctx.Data()
	data.SetData(result)
	return ctx.JSON(data)
}

func FrontendRouteRewriteAdd(ctx echo.Context) error {
	if ctx.Form(`op`) == `routeList` {
		return FrontendRouteList(ctx)
	}
	var err error
	m := official.NewRouteRewrite(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonRouteRewrite, echo.ExcludeFieldName(`updated`))
		if err == nil {
			_, err = m.Add()
			if err != nil {
				goto END
			}
			err = frontend.ResetRouteRewrite()
			if err != nil {
				goto END
			}
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/manager/frontend/route_rewrite`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint64()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCommonRouteRewrite, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

END:
	ctx.Set(`activeURL`, `/manager/frontend/route_rewrite`)
	ctx.Set(`title`, ctx.T(`添加自定义网址`))
	return ctx.Render(`official/manager/frontend/route_rewrite_edit`, err)
}

func FrontendRouteRewriteEdit(ctx echo.Context) error {
	if ctx.Form(`op`) == `routeList` {
		return FrontendRouteList(ctx)
	}
	var err error
	id := ctx.Formx(`id`).Uint64()
	m := official.NewRouteRewrite(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonRouteRewrite, echo.ExcludeFieldName(`created`))
		if err != nil {
			goto END
		}
		err = m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/manager/frontend/route_rewrite`))
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
		echo.StructToForm(ctx, m.OfficialCommonRouteRewrite, ``, echo.LowerCaseFirstLetter)
	}

END:
	ctx.Set(`activeURL`, `/manager/frontend/route_rewrite`)
	ctx.Set(`title`, ctx.T(`编辑自定义网址`))
	return ctx.Render(`official/manager/frontend/route_rewrite_edit`, err)
}

func FrontendRouteRewriteDelete(ctx echo.Context) error {
	ids := ctx.FormxValues(`id`).Uint64(param.IsGreaterThanZeroElement)
	if len(ids) == 0 {
		return ctx.NewError(code.InvalidParameter, `请选择要删除的项`).SetZone(`id`)
	}
	m := official.NewRouteRewrite(ctx)
	var err error
	for _, _v := range ids {
		if err = m.Delete(nil, db.Cond{`id`: _v}); err != nil {
			break
		}
	}
	if err == nil {
		err = frontend.ResetRouteRewrite()
	}
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/manager/frontend/route_rewrite`))
}
