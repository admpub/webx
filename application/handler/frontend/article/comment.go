package article

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webfront/middleware/sessdata"
	modelComment "github.com/coscms/webfront/model/official/comment"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func ArticleCommentAdd(c echo.Context) (err error) {
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

func ArticleCommentReplyList(c echo.Context) error {
	c.SetAuto(true)
	commentID := c.Formx(`commentId`).Uint64()
	c.Request().Form().Set(`pageSize`, `10`)
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

func ArticleCommentList(c echo.Context) (err error) {
	c.SetAuto(true)
	id := c.Formx(`id`).Uint64()
	sn := c.Formx(`sn`).String()
	typ := c.Formx(`type`, `article`).String()
	flat := c.Formx(`flat`).Bool()
	subType := c.Formx(`subtype`).String()
	c.Request().Form().Set(`pageSize`, `10`)
	listx, err := modelComment.QueryCommentList(c, id, sn, typ, subType, flat, ``)
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
