package tool

import (
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

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
		list, err = i18nm.ListByResource(ctx, table)
	}
	ctx.Set(`listData`, list)
	ctx.Set(`langs`, config.FromFile().Language.AllList)
	ctx.Set(`langDefault`, config.FromFile().Language.Default)
	ctx.Set(`title`, ctx.T(`本地化翻译`))
	ctx.SetFunc(`tableTitles`, func() []*echo.KVx[[]string, any] {
		return i18nm.TableTitles.Slice()
	})
	ctx.SetFunc(`tableFields`, func() []string {
		item := i18nm.TableTitles.GetItem(table)
		if item != nil {
			return item.X
		}
		return nil
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
	column := ctx.Form(`column`)
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
	rowID := ctx.Formx(`rowID`).Uint64()
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
	text := ctx.Form(`text`)
	err := i18nm.UpdateColumnTranslation(ctx, table, column, rowID, lang, text)
	if err != nil {
		return err
	}
	data := ctx.Data()
	data.SetInfo(ctx.T(`修改成功`))
	return ctx.JSON(data)
}
