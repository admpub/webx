package advert

import (
	"context"
	"fmt"
	"strings"

	"github.com/webx-top/echo"
)

type (
	ContentRenderer func(Adverter) string
	Adverter        interface {
		GetWidth() uint
		GetHeight() uint
		GetURL() string
		GetContent() string
		GetContype() string
	}
)

var (
	Contype = echo.NewKVData()
	AdMode  = echo.NewKVData()
)

func Render(v Adverter) string {
	item := Contype.GetItem(v.GetContype())
	if item == nil {
		return ``
	}
	if fn, ok := item.Fn()(nil).(ContentRenderer); ok {
		return fn(v)
	}
	return ``
}

func genStyle(p Adverter) string {
	if p == nil {
		return ``
	}
	var styles []string
	if p.GetWidth() > 0 {
		styles = append(styles, fmt.Sprintf(`width:%dpx`, p.GetWidth()))
	}
	if p.GetHeight() > 0 {
		styles = append(styles, fmt.Sprintf(`height:%dpx`, p.GetHeight()))
	}
	style := strings.Join(styles, `;`)
	if len(style) > 0 {
		style = ` style="` + style + `"`
	}
	return style
}

func init() {
	Contype.AddItem(echo.NewKV(`text`, `文字广告`).SetHKV(`description`, `输入广告文字`).SetFn(func(c context.Context) interface{} {
		return ContentRenderer(func(v Adverter) string {
			return `<a href="` + v.GetURL() + `" target="_blank">` + v.GetContent() + `</a>`
		})
	}))
	Contype.AddItem(echo.NewKV(`image`, `图片广告`).SetHKV(`description`, `输入图片文件网址`).SetFn(func(c context.Context) interface{} {
		return ContentRenderer(func(v Adverter) string {
			return `<a href="` + v.GetURL() + `" target="_blank"><img rel="` + v.GetContent() + `" class="previewable" src="` + v.GetContent() + `"` + genStyle(v) + ` /></a>`
		})
	}))
	Contype.AddItem(echo.NewKV(`video`, `视频广告`).SetHKV(`description`, `输入视频文件网址`).SetFn(func(c context.Context) interface{} {
		return ContentRenderer(func(v Adverter) string {
			return `<a href="` + v.GetURL() + `" target="_blank"><video src="` + v.GetContent() + `" controls="controls"` + genStyle(v) + `></video></a>`
		})
	}))
	Contype.AddItem(echo.NewKV(`audio`, `音频广告`).SetHKV(`description`, `输入音频文件网址`).SetFn(func(c context.Context) interface{} {
		return ContentRenderer(func(v Adverter) string {
			return `<a href="` + v.GetURL() + `" target="_blank"><audio src="` + v.GetContent() + `" controls="controls"` + genStyle(v) + `></audio></a>`
		})
	}))

	AdMode.Add(`CPA`, `CPA`)
	AdMode.Add(`CPM`, `CPM`)
	AdMode.Add(`CPC`, `CPC`)
	AdMode.Add(`CPS`, `CPS`)
	AdMode.Add(`CPT`, `CPT`)
}
