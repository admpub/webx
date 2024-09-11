package article

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"

	"github.com/admpub/webx/application/model/official"
	modelArticle "github.com/admpub/webx/application/model/official/article"
	modelComment "github.com/admpub/webx/application/model/official/comment"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
)

func articleFormFilter(options ...formfilter.Options) echo.FormDataFilter {
	options = append(options, formfilter.Exclude(`updated`, `created`, `comments`, `likes`, `hates`, `views`))
	return formfilter.Build(options...)
}

func Index(ctx echo.Context) error {
	if operation := ctx.Form(`operation`); operation == `selectSource` {
		return selectPageSource(ctx)
	}
	sourceID := ctx.Queryx(`sourceId`).String()
	sourceTable := ctx.Queryx(`sourceTable`).String()
	contype := ctx.Queryx(`contype`).String()
	categoryID := ctx.Formx(`categoryId`).Uint()
	m := modelArticle.NewArticle(ctx)
	cond := db.NewCompounds()
	common.SelectPageCond(ctx, cond, `id`, `title%`)
	if len(sourceTable) > 0 {
		if len(sourceID) > 0 {
			cond.AddKV(`source_id`, sourceID)
		}
		cond.AddKV(`source_table`, sourceTable)
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
	list, err := m.ListPage(cond, `-id`)
	ctx.Set(`listData`, list)
	ctx.Set(`sourceId`, sourceID)
	ctx.Set(`sourceTable`, sourceTable)
	ctx.Set(`sourceTables`, modelArticle.Source.Slice())
	ctx.Set(`contypes`, modelArticle.Contype.Slice())
	ctx.SetFunc(`getContypeName`, modelArticle.Contype.Get)
	ctx.SetFunc(`getSourceTableName`, modelArticle.Source.Get)
	return ctx.Render(`official/article/index`, common.Err(ctx, err))
}

func selectPageSource(ctx echo.Context) error {
	sourceTable := ctx.Form(`sourceTable`)
	h := modelArticle.Source.GetSelectPageHandler(sourceTable)
	if h != nil {
		return h(ctx)
	}
	return nil
}

func Add(ctx echo.Context) error {
	sourceID := ctx.Queryx(`sourceId`).String()
	sourceTable := ctx.Queryx(`sourceTable`).String()
	if operation := ctx.Form(`operation`); operation == `selectSource` {
		return selectPageSource(ctx)
	}
	var err error
	m := modelArticle.NewArticle(ctx)
	user := backend.User(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonArticle, articleFormFilter())
		if err == nil {
			m.OwnerId = uint64(user.Id)
			m.OwnerType = `user`
			_, err = m.Add()
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/article/index?sourceId=`) + sourceID + `&sourceTable=` + sourceTable)
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint64()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCommonArticle, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	SetArticleFormData(ctx, sourceID, sourceTable)
	ctx.Set(`activeURL`, `/official/article/index`)
	ctx.Set(`sourceId`, sourceID)
	ctx.Set(`sourceTable`, sourceTable)
	ctx.Set(`contypes`, modelArticle.Contype.Slice())
	return ctx.Render(`official/article/edit`, common.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	sourceID := ctx.Queryx(`sourceId`).String()
	sourceTable := ctx.Queryx(`sourceTable`).String()
	if operation := ctx.Form(`operation`); operation == `selectSource` {
		return selectPageSource(ctx)
	}
	var err error
	id := ctx.Formx(`id`).Uint64()
	m := modelArticle.NewArticle(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonArticle, articleFormFilter())
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				common.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(backend.URLFor(`/official/article/index?sourceId=`) + sourceID + `&sourceTable=` + sourceTable)
			}
		}
	} else if ctx.IsAjax() {
		display := ctx.Query(`display`)
		closeComment := ctx.Query(`closeComment`)
		if len(display) > 0 {
			m.Display = display
			data := ctx.Data()
			err = m.UpdateField(nil, `display`, display, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
		if len(closeComment) > 0 {
			m.CloseComment = closeComment
			data := ctx.Data()
			err = m.UpdateField(nil, `close_comment`, closeComment, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonArticle, ``, echo.LowerCaseFirstLetter)
	}
	SetArticleFormData(ctx, sourceID, sourceTable)
	ctx.Set(`activeURL`, `/official/article/index`)
	ctx.Set(`sourceId`, sourceID)
	ctx.Set(`sourceTable`, sourceTable)
	ctx.Set(`contypes`, modelArticle.Contype.Slice())
	return ctx.Render(`official/article/edit`, common.Err(ctx, err))
}

func SetArticleFormData(ctx echo.Context, sourceID string, sourceTable string) {
	var tags []string
	if len(ctx.Form(`tags`)) > 0 {
		tags = strings.Split(ctx.Form(`tags`), `,`)
	}
	ctx.Set(`tags`, tags)
	ctx.SetFunc(`isCheckedTag`, func(tag string) bool {
		return com.InSlice(tag, tags)
	})
	var tagRows []echo.H
	tagsGetter := modelArticle.Source.GetTagsGetter(sourceTable)
	if tagsGetter != nil {
		tagRows, _ = tagsGetter(ctx, sourceID)
	} else {
		tagsM := official.NewTags(ctx)
		tagList := tagsM.ListByGroup(modelArticle.GroupName, 2000)
		tagRows = make([]echo.H, len(tagList))
		for index, tag := range tagList {
			tagRows[index] = echo.H{`id`: tag.Name, `text`: tag.Name}
		}
	}
	ctx.Set(`tagList`, tagRows)
	ctx.Set(`allowUsers`, modelComment.CommentAllowUsers)
	cateM := official.NewCategory(ctx)
	categoryList := cateM.ListIndent(cateM.ListAllParent(modelArticle.GroupName, 0))
	ctx.Set(`categoryList`, categoryList)
	ctx.Set(`contentHideTags`, modelArticle.ContentHideDetector.Slice())
	ctx.Set(`sourceTableList`, modelArticle.Source.Slice())
}

func Delete(ctx echo.Context) error {
	sourceID := ctx.Queryx(`sourceId`).String()
	sourceTable := ctx.Queryx(`sourceTable`).String()
	id := ctx.Formx(`id`).Uint64()
	m := modelArticle.NewArticle(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/article/index`) + `?sourceId=` + sourceID + `&sourceTable=` + sourceTable)
}
