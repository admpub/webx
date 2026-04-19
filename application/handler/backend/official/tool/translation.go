package tool

import (
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webfront/model/i18nm"
	_ "github.com/coscms/webfront/model/i18nm/listener"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

func registerTranslationResourceTableEditURL() {
	i18nm.RegisterTableTitle(`official_common_article`, echo.T(`文章`), nil, `/official/article/edit?id=`)
	i18nm.RegisterTableTitle(`official_common_category`, echo.T(`分类`), nil, `/official/article/category_edit?id=`)
	i18nm.RegisterTableTitle(`official_common_group`, ``, nil, `/official/customer/group/edit?id=`)
	i18nm.RegisterTableTitle(`official_common_tags`, ``, nil, `/official/tags/edit?id=`)
	i18nm.RegisterTableTitle(`official_common_navigate`, ``, nil, `/manager/navigate/edit?id=`)
	i18nm.RegisterTableTitle(`official_common_route_page`, ``, nil, `/manager/frontend/route_page_edit?id=`)
	i18nm.RegisterTableTitle(`official_customer_role`, ``, nil, `/official/customer/role/edit?id=`)
	i18nm.RegisterTableTitle(`official_customer_level`, ``, nil, `/official/customer/level/edit?id=`)
	i18nm.RegisterTableTitle(`official_customer_group_package`, ``, nil, `/official/customer/group_package/edit?id=`)
	i18nm.RegisterTableTitle(`official_ad_position`, ``, nil, `/official/advert/position_edit?id=`)
	i18nm.RegisterTableTitle(`official_ad_item`, ``, nil, `/official/advert/edit?id=`)
	i18nm.RegisterTableTitle(`official_common_area`, ``, nil, `/tool/area/edit?id=`)
	i18nm.RegisterTableTitle(`official_common_area_country`, ``, nil, `/tool/area/country_edit?id=`)
	i18nm.RegisterTableTitle(`official_common_area_group`, ``, nil, `/tool/area/group_edit?id=`)
}

func translationIndex(ctx echo.Context) error {
	var err error
	table := ctx.Query(`table`)
	if len(table) == 0 {
		item := i18nm.TableTitles.GetItemByIndex(0)
		if item != nil {
			table = item.K
		}
	} else {
		if !i18nm.TableTitles.Has(table) {
			table = ``
			err = ctx.NewError(code.Unsupported, `不支持的表名`).SetZone(`table`)
		}
	}
	var list []echo.H
	if len(table) > 0 {
		common.SetPagingDefaultSize(ctx, 10)
		options := i18nm.ListQuery{
			Table: table,
			RowID: ctx.Queryx(`id`).Uint64(),
			Lang:  ctx.Query(`lang`),
		}
		list, err = i18nm.ListByResource(ctx, options)
	}
	ctx.Set(`listData`, list)
	ctx.Set(`langs`, config.FromFile().Language.AllList)
	ctx.Set(`langDefault`, config.FromFile().Language.Default)
	ctx.Set(`table`, table)
	ctx.Set(`tableTitle`, echo.T(i18nm.TableTitles.Get(table)))
	ctx.Set(`title`, ctx.T(`本地化翻译`))
	ctx.SetFunc(`tableTitles`, func() []*echo.KVx[[]string, any] {
		return i18nm.TableTitles.Slice()
	})
	ctx.SetFunc(`editURL`, func(id any) string {
		item := i18nm.TableTitles.GetItem(table)
		if item != nil {
			url := item.H.String(`editURL`)
			if len(url) > 0 {
				return url + param.AsString(id)
			}
		}
		return ``
	})
	ctx.SetFunc(`tableFields`, func() []string {
		item := i18nm.TableTitles.GetItem(table)
		if item != nil {
			return item.X
		}
		return nil
	})
	ctx.SetFunc(`inputType`, func(field string) string {
		fieldInfo, ok := dbschema.DBI.Fields.Find(table, field)
		if !ok {
			return `text`
		}
		switch fieldInfo.DataType {
		case `text`, `longtext`, `mediumtext`:
			return `textarea`
		case `bigint`, `int`, `tinyint`, `smallint`, `mediumint`:
			return `number`
		case `varchar`:
			if fieldInfo.MaxSize > 255 || field == `summary` || field == `description` {
				return `textarea`
			}
			fallthrough
		default:
			return `text`
		}
	})
	langs := config.FromFile().Language.KVList()
	ctx.SetFunc(`langTitle`, func(lang string) string {
		if len(lang) == 0 {
			lang = config.FromFile().Language.Default
		}
		for _, item := range langs {
			if item.K == lang {
				return item.V
			}
		}
		return lang
	})
	return ctx.Render(`official/tool/translation/index`, common.Err(ctx, err))
}

func translationEdit(ctx echo.Context) error {
	table := ctx.Form(`table`)
	if len(table) == 0 {
		return ctx.NewError(code.InvalidParameter, `缺少表名`).SetZone(`table`)
	}
	if !i18nm.TableTitles.Has(table) {
		return ctx.NewError(code.Unsupported, `不支持的表名`).SetZone(`table`)
	}
	column := ctx.FormAny(`column`, `name`)
	if len(column) == 0 {
		return ctx.NewError(code.InvalidParameter, `缺少列名`).SetZone(`column`)
	}
	fieldInfo, ok := dbschema.DBI.Fields.Find(table, column)
	if !ok {
		return ctx.NewError(code.Unsupported, `不支持的列名`).SetZone(`column`)
	}
	if !fieldInfo.Multilingual {
		return ctx.NewError(code.Unsupported, `该列不支持多语言`).SetZone(`column`)
	}
	rowID := ctx.FormAnyx(`rowID`, `pk`).Uint64()
	if rowID == 0 {
		return ctx.NewError(code.InvalidParameter, `缺少行ID`).SetZone(`rowID`)
	}
	lang := ctx.Form(`lang`)
	if len(lang) == 0 {
		return ctx.NewError(code.InvalidParameter, `缺少语言`).SetZone(`lang`)
	}
	if !com.InSlice(lang, config.FromFile().Language.AllList) {
		return ctx.NewError(code.Unsupported, `不支持的语言`).SetZone(`lang`)
	}
	text := ctx.FormAny(`text`, `value`)
	affected, err := i18nm.UpdateColumnTranslation(ctx, table, column, rowID, lang, text)
	if err != nil {
		return err
	}
	data := ctx.Data()
	if affected == 0 {
		data.SetInfo(ctx.T(`没有修改任何数据`))
	} else {
		data.SetInfo(ctx.T(`修改成功`))
	}
	return ctx.JSON(data)
}
