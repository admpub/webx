package article

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/initialize/frontend"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelComment "github.com/admpub/webx/application/model/official/comment"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func ArticleCommentAdd(c echo.Context) (err error) {
	customer := sessdata.Customer(c)
	data := common.VerifyCaptcha(c, frontend.Name, `code`)
	if common.IsFailureCode(data.GetCode()) {
		return c.JSON(data)
	}
	cpt, err := common.GetCaptchaEngine(c)
	if err != nil {
		return c.JSON(data.SetError(err))
	}
	data.SetData(cpt.MakeData(c, frontend.Name, `code`))
	if customer == nil {
		name := c.Form(`name`)
		pass := c.Form(`password`)
		typi := c.Form(`type`)
		m := modelCustomer.NewCustomer(c)
		err := m.SignIn(name, pass, typi, modelCustomer.GenerateOptionsFromHeader(c)...)
		if err != nil {
			return c.JSON(data.SetError(err))
		}
		customer = m.OfficialCustomer
	}
	typ := c.Formx(`type`, `article`).String()
	tp, ok := modelComment.CommentAllowTypes[typ]
	if !ok {
		return c.NewError(code.InvalidParameter, `不支持评论目标: %v`, typ)
	}
	subType := c.Formx(`subtype`).String()
	id := c.Formx(`id`).Uint64()
	sn := c.Formx(`sn`).String()
	if len(sn) > 0 && tp.SN2ID != nil {
		id, err = tp.SN2ID(c, sn)
		if err != nil {
			return
		}
	}
	cmtM := modelComment.NewComment(c)
	cmtM.TargetType = typ
	cmtM.TargetSubtype = subType
	cmtM.TargetId = id
	cmtM.SetCustomerID(customer.Id)
	cmtM.Content = c.Formx(`content`).String()
	cmtM.Contype = `text`
	cmtM.ReplyCommentId = c.Formx(`replyId`).Uint64()
	_, err = cmtM.Add()
	if err != nil {
		data.SetError(err)
	} else {
		data.SetInfo(c.T(`评论成功`))
	}
	return c.JSON(data)
}

func articleCommentReplyList(c echo.Context, commentID uint64, urlLayout string, pagingVarSuffix ...string) ([]*modelComment.CommentAndExtra, error) {
	if commentID == 0 {
		return nil, c.NewError(code.InvalidParameter, `commentId无效`)
	}
	cmtM := modelComment.NewComment(c)
	err := cmtM.Get(nil, `id`, commentID)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, c.NewError(code.DataNotFound, `评论不存在`)
		}
		return nil, err
	}
	cond := cmtM.ListReplyCond(commentID)
	list := []*modelComment.CommentAndReplyTarget{}
	p, err := common.NewLister(cmtM, &list, func(r db.Result) db.Result {
		return r.OrderBy(`id`).Relation(`ReplyTarget`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
			if modelComment.NeedWithQuoteComment(c) {
				return sel
			}
			return nil
		})
	}, cond.And()).Paging(c, pagingVarSuffix...)
	if err != nil {
		return nil, err
	}
	if len(urlLayout) > 0 {
		p.SetURL(urlLayout)
	}
	return cmtM.WithExtra(list, sessdata.Customer(c), handler.User(c), p)
}

func ArticleCommentReplyList(c echo.Context) error {
	c.SetAuto(true)
	commentID := c.Formx(`commentId`).Uint64()
	c.Request().Form().Set(`pageSize`, `10`)
	listx, err := articleCommentReplyList(c, commentID, ``)
	if err != nil {
		return err
	}
	c.Set(`listComment`, listx)
	c.SetFunc(`commentList`, func() []*modelComment.CommentAndExtra {
		return listx
	})
	c.Set(`isReply`, true)
	c.Request().Form().Set(`_pjax`, `true`)
	if c.Format() == `json` {
		data := c.Data()
		htmlContent := c.Formx(`html`).Bool()
		if htmlContent {
			b, err := c.Fetch(`/article/comment_reply_list`, nil)
			if err != nil {
				return err
			}
			data.SetData(echo.H{
				`html`:       string(b),
				`pagination`: c.Get(`pagination`),
			})
		} else if modelComment.PureJSONCommentList(c) {
			data.SetData(c.Stored())
		}
		return c.JSON(data)
	}
	return c.Render(`/article/comment_reply_list`, nil)
}

func articleCommentList(c echo.Context, articleID uint64, articleSN string, targetType string, subType string, flat bool, urlLayout string, pagingVarSuffix ...string) ([]*modelComment.CommentAndExtra, error) {
	tp, ok := modelComment.CommentAllowTypes[targetType]
	if !ok {
		return nil, c.NewError(code.Unsupported, `不支持评论目标: %v`, targetType).SetZone(`type`)
	}
	var err error
	if len(articleSN) > 0 && tp.SN2ID != nil {
		articleID, err = tp.SN2ID(c, articleSN)
		if err != nil {
			return nil, err
		}
	}
	cmtM := modelComment.NewComment(c)
	cond := cmtM.ListCond(targetType, subType, articleID, flat)
	list := []*modelComment.CommentAndReplyTarget{}
	p, err := common.NewLister(cmtM, &list, func(r db.Result) db.Result {
		return r.OrderBy(`id`).Relation(`ReplyTarget`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
			if modelComment.NeedWithQuoteComment(c) {
				return sel
			}
			return nil
		})
	}, cond.And()).Paging(c, pagingVarSuffix...)
	if err != nil {
		return nil, err
	}
	if len(urlLayout) > 0 {
		p.SetURL(urlLayout)
	}
	var rowNums map[uint64]int
	var replyCommentIndexes map[uint64][]int
	if flat {
		replyCommentIndexes = map[uint64][]int{}
		replyCommentIds := []uint64{}
		for index, row := range list {
			if row.ReplyCommentId > 0 {
				if _, ok := replyCommentIndexes[row.ReplyCommentId]; !ok {
					replyCommentIndexes[row.ReplyCommentId] = []int{}
					replyCommentIds = append(replyCommentIds, row.ReplyCommentId)
				}
				replyCommentIndexes[row.ReplyCommentId] = append(replyCommentIndexes[row.ReplyCommentId], index)
			}
		}
		if len(replyCommentIds) > 0 {
			rowNums, err = cmtM.RowNums(targetType, subType, articleID, replyCommentIds)
			if err != nil {
				return nil, err
			}
		}
	}

	rows, err := cmtM.WithExtra(list, sessdata.Customer(c), handler.User(c), p)
	if err != nil {
		return nil, err
	}
	for id, rowNum := range rowNums {
		for _, index := range replyCommentIndexes[id] {
			rows[index].ReplyFloorNumber = rowNum
		}
	}
	return rows, err
}

func ArticleCommentList(c echo.Context) (err error) {
	c.SetAuto(true)
	id := c.Formx(`id`).Uint64()
	sn := c.Formx(`sn`).String()
	typ := c.Formx(`type`, `article`).String()
	flat := c.Formx(`flat`).Bool()
	subType := c.Formx(`subtype`).String()
	c.Request().Form().Set(`pageSize`, `10`)
	listx, err := articleCommentList(c, id, sn, typ, subType, flat, ``)
	if err != nil {
		return err
	}
	c.Set(`listComment`, listx)
	c.Set(`flat`, flat)
	c.SetFunc(`commentList`, func() []*modelComment.CommentAndExtra {
		return listx
	})
	c.Set(`isReply`, false)
	c.Request().Form().Set(`_pjax`, `true`)
	if c.Format() == `json` {
		data := c.Data()
		htmlContent := c.Formx(`html`).Bool()
		if htmlContent {
			b, err := c.Fetch(`/article/comment_list`, nil)
			if err != nil {
				return err
			}
			data.SetData(echo.H{
				`html`:       string(b),
				`pagination`: c.Get(`pagination`),
			})
		} else if modelComment.PureJSONCommentList(c) {
			data.SetData(c.Stored())
		}
		return c.JSON(data)
	}
	return c.Render(`/article/comment_list`, nil)
}
