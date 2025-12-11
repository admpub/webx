package manager

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/admpub/webx/application/handler/backend/official/page"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/coscms/webfront/model/official"
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
	return ctx.Render(`official/manager/frontend/route_page`, common.Err(ctx, err))
}

func FrontendRoutePageAdd(ctx echo.Context) error {
	var err error
	m := official.NewRoutePage(ctx)
	templateFiles := page.ListTemplateFileNames(`route_page`)
	if ctx.IsGet() {
		id := ctx.Formx(`copyId`).Uint64()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				m.Id = 0
				i18nm.SetModelTranslationsToForm(ctx, m.OfficialCommonRoutePage, id)
			}
		}

	}
	form := formbuilder.New(ctx,
		m.OfficialCommonRoutePage,
		formbuilder.ConfigFile(`official/manager/frontend/route_page_edit`),
		formbuilder.AllowedNames(
			`pageType`, `method`, `route`, `pageVars`, `templateEnabled`, `templateFile`, ``, `disabled`, `name`, `pageContent`,
		),
	)
	form.OnPost(func() error {
		m.Method = strings.Join(ctx.FormxValues(`method[]`), `,`)
		_, err := m.Add()
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialCommonRoutePage, m.Id, i18nm.OptionContentType(`pageContent`, m.PageType))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/manager/frontend/route_page`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	nameField := form.MultilingualField(config.FromFile().Language.Default, `name`, `name`)
	nameField.AddTag(`required`)

	ctx.Set(`activeURL`, `/manager/frontend/route_page`)
	ctx.Set(`title`, ctx.T(`添加自定义页面`))
	ctx.Set(`typeList`, official.RoutePageTypes.Slice())
	ctx.Set(`methodList`, echo.Methods())
	ctx.SetFunc(`methodChecked`, func(method string, methods []string) bool {
		return com.InSlice(method, methods)
	})
	ctx.Set(`templateFiles`, templateFiles)
	return ctx.Render(`official/manager/frontend/route_page_edit`, err)
}

func FrontendRoutePageEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint64()
	m := official.NewRoutePage(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}

	if ctx.IsGet() {
		if ctx.IsAjax() {
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
		}
		i18nm.SetModelTranslationsToForm(ctx, m.OfficialCommonRoutePage, id)
	}
	form := formbuilder.New(ctx,
		m.OfficialCommonRoutePage,
		formbuilder.ConfigFile(`official/manager/frontend/route_page_edit`),
		formbuilder.AllowedNames(
			`pageType`, `method`, `route`, `pageVars`, `templateEnabled`, `templateFile`, ``, `disabled`, `name`, `pageContent`,
		),
	)
	form.OnPost(func() error {
		m.Method = strings.Join(ctx.FormxValues(`method[]`), `,`)
		err := m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialCommonRoutePage, m.Id, i18nm.OptionContentType(`pageContent`, m.PageType))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/manager/frontend/route_page`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	nameField := form.MultilingualField(config.FromFile().Language.Default, `name`, `name`)
	nameField.AddTag(`required`)

	ctx.Set(`activeURL`, `/manager/frontend/route_page`)
	ctx.Set(`title`, ctx.T(`编辑自定义页面`))
	ctx.Set(`typeList`, official.RoutePageTypes.Slice())
	ctx.Set(`methodList`, echo.Methods())
	ctx.SetFunc(`methodChecked`, func(method string, methods []string) bool {
		return com.InSlice(method, methods)
	})
	templateFiles := page.ListTemplateFileNames(`route_page`)
	ctx.Set(`templateFiles`, templateFiles)
	return ctx.Render(`official/manager/frontend/route_page_edit`, err)
}

func FrontendRoutePageDelete(ctx echo.Context) error {
	ids := ctx.FormxValues(`id`).Unique().Uint64(param.IsGreaterThanZeroElement)
	if len(ids) == 0 {
		return ctx.NewError(code.InvalidParameter, `请选择要删除的项`).SetZone(`id`)
	}
	m := official.NewRoutePage(ctx)
	var err error
	for _, _v := range ids {
		if err = m.Delete(nil, db.Cond{`id`: _v}); err != nil {
			break
		}
	}
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/manager/frontend/route_page`))
}
