package article

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/coscms/webfront/model/official"
)

// ClickFlow 表态
func ClickFlow(c echo.Context, typ string, targetType string, canCancel ...bool) error {
	customer := sessdata.Customer(c)
	data := c.Data()
	if customer == nil {
		data.SetError(nerrors.ErrUserNotLoggedIn)
		return c.JSON(data)
	}
	var id interface{}
	inputID := c.Param(`id`)
	if len(inputID) == 0 {
		inputID = c.Form(`id`)
	}
	var targetID uint64
	if len(inputID) > 0 {
		targetID = param.AsUint64(inputID)
		if targetID == 0 {
			data.SetInfo(c.T(`id无效`), 0)
			return c.JSON(data)
		}
		id = targetID
	} else {
		targetSN := c.Param(`sn`)
		if len(targetSN) == 0 {
			targetSN = c.Formx(`sn`).String()
		}
		if len(targetSN) == 0 {
			data.SetInfo(c.T(`id无效`), 0)
			return c.JSON(data)
		}
		id = targetSN
	}
	target, ok := official.ClickFlowTargets[targetType]
	if !ok {
		return c.NewError(code.Unsupported, `不支持的目标类型: %s`, targetType)
	}
	after, infoGetter, err := target.Do(c, id)
	if err != nil {
		data.SetError(err)
		return c.JSON(data)
	}
	if infoGetter != nil {
		info := infoGetter()
		if info.ID > 0 {
			targetID = info.ID
		}
		if info.AuthorID > 0 && customer.Id == info.AuthorID {
			return c.NewError(code.Unsupported, `操作失败：您不能对自己发布的内容进行此项操作`).SetZone(`type`)
		}
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
		err = c.NewError(code.InvalidParameter, `类型无效: %s`, typ)
	}
	if err == nil {
		if after != nil {
			err = after(typ)
		}
	} else {
		if len(canCancel) > 0 && canCancel[0] && echo.IsErrorCode(err, code.RepeatOperation) {
			err = clickFlowM.DelByTargetOwner(targetType, targetID, clickFlowM.OwnerId, clickFlowM.OwnerType)
			if err == nil {
				if after != nil {
					err = after(typ, true)
				}
				data.SetData(echo.H{`cancel`: true}, code.Success.Int())
			}
		} else if echo.IsErrorCode(err, code.DataAlreadyExists) {
			data.SetData(echo.H{`hasOther`: true})
		}
	}
	if err != nil {
		data.SetError(err)
	} else {
		data.SetInfo(c.T(`表态成功`), code.Success.Int())
	}
	return c.JSON(data)
}

// ArticleLike 文章:喜欢
func ArticleLike(c echo.Context) error {
	return ClickFlow(c, `like`, `article`)
}

// ArticleHate 文章:不喜欢
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

// ArticleLike 文章:收藏
func ArticleCollect(c echo.Context) error {
	return Collect(c, `article`, true)
}

// Collect 收藏
func Collect(c echo.Context, targetType string, canCancel ...bool) error {
	customer := sessdata.Customer(c)
	data := c.Data()
	if customer == nil {
		data.SetError(nerrors.ErrUserNotLoggedIn)
		return c.JSON(data)
	}
	inputID := c.Param(`id`)
	if len(inputID) == 0 {
		inputID = c.Form(`id`)
	}
	var id interface{}
	var targetID uint64
	if len(inputID) > 0 {
		targetID = param.AsUint64(inputID)
		if targetID == 0 {
			data.SetInfo(c.T(`id无效`), 0)
			return c.JSON(data)
		}
		id = targetID
	} else {
		targetSN := c.Param(`sn`)
		if len(targetSN) == 0 {
			targetSN = c.Formx(`sn`).String()
		}
		if len(targetSN) == 0 {
			data.SetInfo(c.T(`id无效`), 0)
			return c.JSON(data)
		}
		id = targetSN
	}
	target, ok := official.CollectionTargets[targetType]
	if !ok {
		return c.E(`不支持的目标类型: %s`, targetType)
	}
	after, infoGetter, err := target.Do(c, id)
	if err != nil {
		data.SetError(err)
		return c.JSON(data)
	}
	collectionM := official.NewCollection(c)
	collectionM.TargetType = targetType
	if infoGetter != nil {
		info := infoGetter()
		if info.ID > 0 {
			targetID = info.ID
		}
		collectionM.Title = info.Title
	}
	collectionM.TargetId = targetID
	collectionM.CustomerId = customer.Id
	_, err = collectionM.Add()
	if err != nil {
		if len(canCancel) > 0 && canCancel[0] && echo.IsErrorCode(err, code.RepeatOperation) {
			err = collectionM.DelByTargetOwner(targetType, targetID, collectionM.CustomerId)
			if err == nil {
				if after != nil {
					err = after(true)
				}
				data.SetData(echo.H{`cancel`: true}, code.Success.Int())
			}
		}
	} else {
		if after != nil {
			err = after()
		}
	}
	if err != nil {
		data.SetError(err)
	} else {
		data.SetInfo(c.T(`收藏成功`), code.Success.Int())
	}
	return c.JSON(data)
}
