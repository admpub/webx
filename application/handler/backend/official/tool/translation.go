package tool

import (
	"github.com/coscms/webfront/model/i18nm"
	"github.com/webx-top/echo"
)

func translationIndex(ctx echo.Context) error {
	i18nm.List
	ctx.Set(`title`, ctx.T(`本地化翻译`))
	return ctx.Render(`official/tool/translation/index`, nil)
}
