package tool

import (
	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webfront/library/search/segment"

	//_ "github.com/coscms/webfront/library/search/segment/gojieba"
	//_ "github.com/coscms/webfront/library/search/segment/jiebago"
	_ "github.com/coscms/webfront/library/search/segment/sego"
)

var SegmentMode = echo.NewKVData()

func init() {
	SegmentMode.Add(`all`, `完整模式`)
	SegmentMode.Add(`new`, `新词识别`)
	SegmentMode.Add(`search`, `搜索引擎模式`)
	SegmentMode.Add(`tag`, `词性标注`)
	SegmentMode.Add(`keywords`, `关键词提取`)
}

func Segment(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		keywords := ctx.Form(`keywords`)
		if len(keywords) == 0 {
			return ctx.JSON(ctx.Data().SetInfo(ctx.T(`关键词不能为空`), code.InvalidParameter.Int()))
		}
		mode := ctx.Form(`mode`)
		if !SegmentMode.Has(mode) {
			return ctx.JSON(ctx.Data().SetInfo(ctx.T(`模式不正确`), code.InvalidParameter.Int()))
		}
		splitedWords := segment.SplitWordsBy(com.Str2bytes(keywords), mode)
		return ctx.JSON(ctx.Data().SetData(splitedWords))
	}
	ctx.Set(`modes`, SegmentMode.Slice())
	ctx.Set(`isInitialized`, segment.IsInitialized())
	return ctx.Render(`official/tool/segment/index`, common.Err(ctx, err))
}
