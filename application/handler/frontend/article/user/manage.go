package user

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	hanlderArticle "github.com/admpub/webx/application/handler/backend/official/article"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelArticle "github.com/admpub/webx/application/model/official/article"
)

func ListByCustomer(ctx echo.Context, customer *dbschema.OfficialCustomer) error {
	m := modelArticle.NewArticle(ctx)
	cond := &db.Compounds{
		db.Cond{`owner_id`: customer.Id},
		db.Cond{`owner_type`: `customer`},
	}
	listFilterer(ctx, cond, m)
	cond.AddKV(`display`, `Y`)
	sorts := common.Sorts(ctx, `official_common_article`, `-id`)
	list, err := m.ListPage(cond, sorts...)
	ctx.Set(`list`, list)
	return err
}

func listFilterer(ctx echo.Context, cond *db.Compounds, m *modelArticle.Article) {
	contype := ctx.Queryx(`contype`).String()
	categoryID := ctx.Formx(`categoryId`).Uint()
	title := ctx.Formx(`q`).String()
	if len(title) > 0 {
		cond.AddKV(`title`, title)
	}
	if categoryID > 0 {
		cond.Add(db.Or(
			db.Cond{`category1`: categoryID},
			db.Cond{`category2`: categoryID},
			db.Cond{`category3`: categoryID},
			db.Cond{`category_id`: categoryID},
		))
	}
	if len(contype) > 0 {
		cond.AddKV(`contype`, contype)
	}
	tag := ctx.Formx(`tag`).String()
	if len(tag) > 0 {
		cond.Add(m.TagCond(tag))
	}
	size := ctx.Formx(`size`).Int()
	if size < 1 {
		ctx.Request().Form().Set(`size`, `20`)
	}
}

// List 用户文章列表
func List(ctx echo.Context) error {
	customer := sessdata.Customer(ctx)
	var err error
	m := modelArticle.NewArticle(ctx)
	cond := &db.Compounds{
		db.Cond{`owner_id`: customer.Id},
		db.Cond{`owner_type`: `customer`},
	}
	listFilterer(ctx, cond, m)
	sorts := common.Sorts(ctx, `official_common_article`, `-id`)
	list, err := m.ListPage(cond, sorts...)
	ctx.Set(`list`, list)
	ctx.SetFunc(`getContypeName`, modelArticle.Contype.Get)
	return ctx.Render(`article/user/list`, handler.Err(ctx, err))
}

func applyFormData(ctx echo.Context, m *dbschema.OfficialCommonArticle) {
	m.CategoryId = ctx.Formx(`categoryId`).Uint()
	m.Title = ctx.Formx(`title`).String()
	m.Image = ctx.Formx(`image`).String()
	m.ImageOriginal = ctx.Formx(`imageOriginal`).String()
	m.Summary = ctx.Formx(`summary`).String()
	m.Summary = ctx.Formx(`summary`).String()
	m.Content = ctx.Formx(`content`).String()
	m.Contype = ctx.Formx(`contype`).String()
	m.Display = `N`
	m.Tags = ctx.Formx(`tags`).String()
}

// Create 创建文章
func Create(ctx echo.Context) error {
	sourceID := ``
	sourceTable := ``
	customer := sessdata.Customer(ctx)
	var err error
	m := modelArticle.NewArticle(ctx)
	m.DisallowCreateTags = true
	if ctx.IsPost() {
		m.OwnerId = customer.Id
		m.OwnerType = `customer`
		applyFormData(ctx, m.OfficialCommonArticle)
		_, err = m.Add()
		if err != nil {
			goto END
		}
		return ctx.Redirect(sessdata.URLFor(`/user/articles`))
	}

END:
	hanlderArticle.SetArticleFormData(ctx, sourceID, sourceTable)
	ctx.Set(`activeURL`, `/user/articles`)
	ctx.Set(`sourceId`, sourceID)
	ctx.Set(`sourceTable`, sourceTable)
	ctx.Set(`contypes`, modelArticle.Contype.Slice())
	ctx.Set(`title`, ctx.T(`投稿`))
	ctx.Set(`isEdit`, false)
	return ctx.Render(`article/user/edit`, handler.Err(ctx, err))
}

// Edit 修改文章
func Edit(ctx echo.Context) error {
	sourceID := ``
	sourceTable := ``
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.NewError(code.InvalidParameter, `参数“%s”值无效`, `id`)
	}
	customer := sessdata.Customer(ctx)
	m := modelArticle.NewArticle(ctx)
	m.DisallowCreateTags = true
	err := m.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `文章不存在`)
		}
		return err
	}
	if m.OwnerType != `customer` || m.OwnerId != customer.Id {
		return ctx.NewError(code.NonPrivileged, `越权操作！您没有权限修改此数据`)
	}
	if ctx.IsPost() {
		applyFormData(ctx, m.OfficialCommonArticle)
		err = m.Edit(nil, `id`, m.Id)
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`修改成功`))
		return ctx.Redirect(sessdata.URLFor(`/user/articles`))
	}
	echo.StructToForm(ctx, m.OfficialCommonArticle, ``, echo.LowerCaseFirstLetter)

END:
	hanlderArticle.SetArticleFormData(ctx, sourceID, sourceTable)
	ctx.Set(`activeURL`, `/user/articles`)
	ctx.Set(`sourceId`, sourceID)
	ctx.Set(`sourceTable`, sourceTable)
	ctx.Set(`contypes`, modelArticle.Contype.Slice())
	ctx.Set(`title`, ctx.T(`修改文章`))
	ctx.Set(`isEdit`, true)
	return ctx.Render(`article/user/edit`, common.Err(ctx, err))
}

// Delete 用户删除文章
func Delete(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.NewError(code.InvalidParameter, `参数“%s”值无效`, `id`)
	}
	customer := sessdata.Customer(ctx)
	m := modelArticle.NewArticle(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `文章不存在`)
		}
		return err
	}
	if m.OwnerType != `customer` || m.OwnerId != customer.Id {
		return ctx.NewError(code.NonPrivileged, `越权操作！您没有权限删除此数据`)
	}
	return ctx.Redirect(sessdata.URLFor(`/user/articles`))
}
