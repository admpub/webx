package comment

import (
	"errors"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"
	modelArticle "github.com/admpub/webx/application/model/official/article"
	modelComment "github.com/admpub/webx/application/model/official/comment"
)

func Index(ctx echo.Context) error {
	q := ctx.Form(`q`)
	display := ctx.Form(`display`)
	m := modelComment.NewComment(ctx)
	cond := []db.Compound{}
	if len(display) > 0 {
		cond = append(cond, db.Cond{`display`: display})
	}
	if len(q) > 0 {
		cond = append(cond, db.Cond{`content`: db.Like(`%` + q + `%`)})
	}
	ctx.Request().Form().Set(`pageSize`, `10`)
	list := []*modelComment.CommentAndReplyTarget{}
	p, err := handler.PagingWithLister(ctx, handler.NewLister(m, &list, func(r db.Result) db.Result {
		return r.OrderBy(`-id`).Relation(`ReplyTarget`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
			if modelComment.NeedWithQuoteComment(ctx) {
				return sel
			}
			return nil
		})
	}, db.And(cond...)))
	ret := handler.Err(ctx, err)
	listx, err := m.WithExtra(list, nil, handler.User(ctx), p)
	if err != nil {
		return err
	}
	ctx.Set(`data`, nil)
	ctx.Set(`listData`, listx)
	ctx.Set(`contypes`, modelArticle.Contype.Slice())
	if ctx.Formx(`partial`).Bool() {
		return ctx.Render(`official/comment/list_partial`, ret)
	}
	return ctx.Render(`official/comment/list`, ret)
}

func List(ctx echo.Context) error {
	q := ctx.Form(`q`)
	display := ctx.Form(`display`)
	targetType := ctx.Form(`targetType`, `article`)
	targetID := ctx.Formx(`targetId`).Uint64()
	m := modelComment.NewComment(ctx)
	cond := []db.Compound{
		//db.Cond{`reply_comment_id`: 0},
	}
	if len(targetType) > 0 && targetID > 0 {
		cond = append(cond, db.Cond{`target_type`: targetType})
		cond = append(cond, db.Cond{`target_id`: targetID})
	}
	if len(display) > 0 {
		cond = append(cond, db.Cond{`display`: display})
	}
	if len(q) > 0 {
		cond = append(cond, db.Cond{`content`: db.Like(`%` + q + `%`)})
	}
	tp, ok := modelComment.CommentAllowTypes[targetType]
	if !ok {
		return ctx.E(`不支持评论目标: %v`, targetType)
	}
	if tp.GetTarget != nil {
		data, breadcrumb, detailURL, err := tp.GetTarget(ctx, targetID)
		if err != nil {
			return err
		}
		ctx.Set(`data`, data)
		ctx.Set(`breadcrumb`, breadcrumb)
		ctx.Set(`targetDetailURL`, detailURL)
	}
	ctx.Request().Form().Set(`pageSize`, `10`)
	list := []*modelComment.CommentAndReplyTarget{}
	p, err := handler.PagingWithLister(ctx, handler.NewLister(m, &list, func(r db.Result) db.Result {
		return r.OrderBy(`-id`).Relation(`ReplyTarget`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
			if modelComment.NeedWithQuoteComment(ctx) {
				return sel
			}
			return nil
		})
	}, db.And(cond...)))
	ret := handler.Err(ctx, err)
	listx, err := m.WithExtra(list, nil, handler.User(ctx), p)
	if err != nil {
		return err
	}
	ctx.Set(`listData`, listx)
	ctx.Set(`contypes`, modelArticle.Contype.Slice())
	if ctx.Formx(`partial`).Bool() {
		return ctx.Render(`official/comment/list_partial`, ret)
	}
	return ctx.Render(`official/comment/list`, ret)
}

func Add(ctx echo.Context) error {
	targetType := ctx.Form(`targetType`, `article`)
	targetID := ctx.Formx(`targetId`).Uint64()
	if targetID == 0 {
		return ctx.E(`targetId无效`)
	}
	articleM := modelArticle.NewArticle(ctx)
	err := articleM.Get(nil, `id`, targetID)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = errors.New(`文章不存在`)
		}
		return err
	}
	m := modelComment.NewComment(ctx)
	user := handler.User(ctx)
	if ctx.IsPost() {
		m.TargetType = targetType
		m.Contype = ctx.Formx(`contype`).String()
		m.TargetId = targetID
		m.CopyFromUser(user)
		m.Content = ctx.Formx(`content`).String()
		m.ReplyCommentId = ctx.Formx(`replyId`).Uint64()
		_, err = m.Add()
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/official/article/comment/index`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint64()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.OfficialCommonComment, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	ctx.Set(`activeURL`, `/official/article/comment/index`)
	ctx.Set(`allowUsers`, modelComment.CommentAllowUsers)
	return ctx.Render(`official/comment/edit`, handler.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint64()
	m := modelComment.NewComment(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCommonComment, echo.ExcludeFieldName(`updated`, `created`, `replies`, `likes`, `hates`))
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/official/article/comment/index`))
			}
		}
	} else if ctx.IsAjax() {
		display := ctx.Query(`display`)
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
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCommonComment, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	ctx.Set(`activeURL`, `/official/article/comment/index`)
	ctx.Set(`allowUsers`, modelComment.CommentAllowUsers)
	return ctx.Render(`official/comment/edit`, handler.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := modelComment.NewComment(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/official/article/comment/index`))
}
