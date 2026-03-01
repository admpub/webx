package article

import (
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webfront/library/xkv"
	"github.com/coscms/webfront/middleware/sessdata"
	modelComment "github.com/coscms/webfront/model/official/comment"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func articleCommentAdd(c echo.Context) error {
	return CommentAdd(c, c.Formx(`type`, `article`).String(), c.Formx(`subtype`).String())
}

func articleCommentReplyList(c echo.Context) error {
	return CommentReplyList(c, `/article/comment_reply_list`)
}

func articleCommentList(c echo.Context) error {
	return CommentList(c, c.Formx(`type`, `article`).String(), c.Formx(`subtype`).String(), `/article/comment_list`)
}

func SetCommentData(c echo.Context, commentURLLayout string, targetType string, targetSubtype string, id uint64) func(disabledMsg error, needReview bool) {
	flat := true
	if sessdata.IsAdmin(c) {
		flat = c.Formx(`flat`, `1`).Bool()
	}
	v, _ := xkv.GetValue(c, `ARTICLE_COMMENT_FLAT`, `1`, `文章评论是否扁平化显示`)
	if len(v) > 0 {
		flat = param.AsBool(v)
	}
	c.Set(`flat`, flat)
	c.SetFunc(`commentList`, func() []*modelComment.CommentAndExtra {
		c.Request().Form().Set(`_pjax`, `true`)
		pagination.SetForceSize(c, 10)
		commentList, _ := modelComment.QueryCommentList(c, id, ``, targetType, targetSubtype, flat, commentURLLayout, `Comment`)
		//c.Request().Form().Del(`_pjax`)
		return commentList
	})

	return func(disabledMsg error, needReview bool) {
		// 是否允许评论
		c.Set(`disabledCommentMessage`, disabledMsg)
		c.Set(`needReviewComment`, needReview)
	}
}

func CommentAdd(c echo.Context, commentType string, subType string) (err error) {
	customer := sessdata.Customer(c)
	data := captchabiz.VerifyCaptcha(c, httpserver.KindFrontend, `code`)
	if nerrors.IsFailureCode(data.GetCode()) {
		return c.JSON(data)
	}
	cpt, err := captchabiz.GetCaptchaEngine(c)
	if err != nil {
		return c.JSON(data.SetError(err))
	}
	data.SetData(cpt.MakeData(c, httpserver.KindFrontend, `code`))
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
	tp, ok := modelComment.CommentAllowTypes[commentType]
	if !ok {
		return c.NewError(code.InvalidParameter, `不支持评论目标: %v`, commentType)
	}
	id := c.Formx(`id`).Uint64()
	sn := c.Formx(`sn`).String()
	if len(sn) > 0 && tp.SN2ID != nil {
		id, err = tp.SN2ID(c, sn)
		if err != nil {
			return
		}
	}
	cmtM := modelComment.NewComment(c)
	cmtM.TargetType = commentType
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

func CommentReplyList(c echo.Context, templ string) error {
	c.SetAuto(true)
	commentID := c.Formx(`commentId`).Uint64()
	pagination.SetForceSize(c, 10)
	listx, err := modelComment.QueryCommentReplyList(c, commentID, ``)
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
			b, err := c.Fetch(templ, nil)
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
	return c.Render(templ, nil)
}

func CommentList(c echo.Context, commentType string, subType string, templ string) error {
	c.SetAuto(true)
	id := c.Formx(`id`).Uint64()
	sn := c.Formx(`sn`).String()
	flat := c.Formx(`flat`).Bool()
	pagination.SetForceSize(c, 10)
	listx, err := modelComment.QueryCommentList(c, id, sn, commentType, subType, flat, ``)
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
			b, err := c.Fetch(templ, nil)
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
	return c.Render(templ, nil)
}
