package official

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/null"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/middleware/sessdata"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

var NavigateLinkType = echo.NewKVData()

func init() {
	NavigateLinkType.Add(`custom`, `自定义链接`)
	item := echo.NewKV(`article-category`, `文章分类`)
	item.SetFn(func(c context.Context) interface{} {
		ctx := c.(echo.Context)
		m := NewCategory(ctx)
		categories := m.ListAllParentBy(`article`, 0, 2, db.Cond{`show_on_menu`: `Y`})
		var (
			list     []*NavigateExt
			children = map[uint][]*NavigateExt{}
			ext      = ctx.DefaultExtension()
		)
		for _, category := range categories {
			navExt := &NavigateExt{
				OfficialCommonNavigate: &dbschema.OfficialCommonNavigate{
					Id:       category.Id,
					ParentId: category.ParentId,
					Title:    category.Name,
					Url:      sessdata.URLFor(`/articles` + ext + `?categoryId=` + fmt.Sprint(category.Id)),
				},
				Extra: echo.H{
					`object`: category,
				},
				isActive: GenActiveDetector(`categoryId`, category.Id),
			}
			navExt.SetContext(ctx)
			if category.ParentId < 1 {
				list = append(list, navExt)
				continue
			}
			if _, ok := children[category.ParentId]; !ok {
				children[category.ParentId] = []*NavigateExt{}
			}
			children[category.ParentId] = append(children[category.ParentId], navExt)
		}
		FillChildren(&list, children)
		return list
	})
	NavigateLinkType.AddItem(item)
}

func FillChildren(list *[]*NavigateExt, children map[uint][]*NavigateExt) {
	if len(children) == 0 {
		return
	}
	for _, item := range *list {
		subItems, ok := children[item.Id]
		if !ok {
			continue
		}
		item.Children = &subItems
		delete(children, item.Id)
		FillChildren(item.Children, children)
	}
}

func NewNavigateExt(nav *dbschema.OfficialCommonNavigate) *NavigateExt {
	return &NavigateExt{OfficialCommonNavigate: nav}
}

func GenActiveDetector(categoryKey string, categoryID uint) func(ctx echo.Context) bool {
	return func(ctx echo.Context) bool {
		return ctx.Formx(categoryKey).Uint() == categoryID
	}
}

type NavigateExt struct {
	*dbschema.OfficialCommonNavigate
	isInside null.Bool // 是否是内部链接
	isActive func(echo.Context) bool
	Children *[]*NavigateExt
	Extra    echo.H
}

func (f *NavigateExt) IsInside() bool {
	if !f.isInside.Valid {
		f.isInside.Bool = strings.HasPrefix(f.Url, `/`)
		f.isInside.Valid = true
	}
	return f.isInside.Bool
}

func (f *NavigateExt) URL() string {
	if len(f.Url) == 0 {
		return ``
	}
	if f.IsInside() {
		return sessdata.URLFor(f.Url)
	}
	return f.Url
}

func (f *NavigateExt) IsValidURL() bool {
	if len(f.Url) == 0 || strings.HasPrefix(f.Url, `#`) {
		return false
	}
	return true
}

func (f *NavigateExt) IsActive() bool {
	if f.isActive != nil {
		return f.isActive(f.Context())
	}
	currentPath := f.Context().Request().URL().Path()
	if f.IsInside() && currentPath == f.Url {
		return true
	}
	if len(f.Ident) > 0 {
		if strings.HasPrefix(f.Ident, `regexp:`) {
			expr := strings.TrimPrefix(f.Ident, `regexp:`)
			re, err := regexp.Compile(expr)
			if err != nil {
				log.Error(expr+`: `, err)
				return false
			}
			return re.MatchString(currentPath)
		}
		return strings.HasSuffix(currentPath, f.Ident)
	}
	return false
}

func (f *NavigateExt) SetActiveDetector(fn func(echo.Context) bool) *NavigateExt {
	f.isActive = fn
	return f
}

func (f *NavigateExt) SetExtra(extra echo.H) *NavigateExt {
	f.Extra = extra
	return f
}

func (f *NavigateExt) SetExtraKV(k string, v interface{}) *NavigateExt {
	if f.Extra == nil {
		f.Extra = echo.H{}
	}
	f.Extra.Set(k, v)
	return f
}

func (f *NavigateExt) HasChildren() bool {
	f.getChildren()
	return len(*f.Children) > 0
}

func (f *NavigateExt) FetchChildren(forces ...bool) []*NavigateExt {
	f.getChildren(forces...)
	return *f.Children
}

func (f *NavigateExt) ClearChildren() *NavigateExt {
	f.Children = nil
	return f
}

func (f *NavigateExt) getChildren(forces ...bool) *NavigateExt {
	var force bool
	if len(forces) > 0 {
		force = forces[0]
	}
	if !force && f.Children != nil {
		return f
	}
	defer func() {
		if f.Children == nil {
			f.Children = &[]*NavigateExt{}
		}
	}()
	if len(f.LinkType) == 0 {
		return f
	}
	item := NavigateLinkType.GetItem(f.LinkType)
	if item == nil {
		return f
	}
	if item.Fn() == nil {
		return f
	}
	if list, ok := item.Fn()(f.Context()).([]*NavigateExt); ok {
		f.Children = &list
	}
	return f
}
