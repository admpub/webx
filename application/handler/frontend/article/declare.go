package article

import (
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/coscms/webfront/model/official"
)

// ClickFlow 表态
func ClickFlow(c echo.Context, typ string, targetType string) error {
	customer := sessdata.Customer(c)
	data := c.Data()
	if customer == nil {
		data.SetError(nerrors.ErrUserNotLoggedIn)
		return c.JSON(data)
	}
	var id interface{}
	targetID := c.Formx(`id`).Uint64()
	targetSN := c.Formx(`sn`).String()
	if targetID == 0 {
		if len(targetSN) > 0 {
			id = targetSN
		} else {
			data.SetInfo(c.T(`id无效`), 0)
			return c.JSON(data)
		}
	} else {
		id = targetID
	}
	target, ok := official.ClickFlowTargets[targetType]
	if !ok {
		return c.E(`不支持的目标类型: %s`, targetType)
	}
	after, idGetter, err := target.Do(c, id)
	if err != nil {
		data.SetError(err)
		return c.JSON(data)
	}
	if idGetter != nil {
		targetID = idGetter()
	}
	clickFlowM := official.NewClickFlow(c)
	clickFlowM.TargetType = targetType
	clickFlowM.TargetId = targetID
	clickFlowM.OwnerId = customer.Id
	clickFlowM.OwnerType = `customer`
	clickFlowM.Type = typ
	switch typ {
	case `like`:
		_, err = clickFlowM.Add()
	case `hate`:
		_, err = clickFlowM.Add()
	default:
		err = c.E(`类型无效`)
	}
	if err == nil {
		if after != nil {
			err = after(typ)
		}
	}
	if err != nil {
		data.SetError(err)
	} else {
		data.SetInfo(c.T(`表态成功`))
	}
	return c.JSON(data)
}

// ArticleLike 新闻:喜欢
func ArticleLike(c echo.Context) error {
	return ClickFlow(c, `like`, `article`)
}

// ArticleHate 新闻:不喜欢
func ArticleHate(c echo.Context) error {
	return ClickFlow(c, `hate`, `article`)
}

// CommentLike 评论:喜欢
func CommentLike(c echo.Context) error {
	return ClickFlow(c, `like`, `comment`)
}

// CommentHate 评论:不喜欢
func CommentHate(c echo.Context) error {
	return ClickFlow(c, `hate`, `comment`)
}
