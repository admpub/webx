package userhome

import (
	"strings"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/dashboard"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	registryUserhome "github.com/coscms/webfront/registry/userhome"
	"github.com/coscms/webfront/transform/transformCustomer"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func Index(ctx echo.Context) error {
	uidParam := ctx.Paramx(`customerId`)
	operate := ctx.Param(`operate`)
	cond := db.NewCompounds()
	name := uidParam.String()
	if strings.HasPrefix(name, `@`) { // @username
		name = strings.TrimPrefix(name, `@`)
		if len(name) == 0 {
			return ctx.NewError(code.InvalidParameter, `非法参数`)
		}
		cond.AddKV(`name`, name)
	} else {
		customerID := uidParam.Uint64()
		if customerID < 1 {
			return ctx.NewError(code.InvalidParameter, `非法参数`)
		}
		cond.AddKV(`id`, customerID)
	}
	m := modelCustomer.NewCustomer(ctx)
	detail, err := m.GetDetail(cond.And())
	if err != nil {
		if err == db.ErrNoMoreRows {
			return ctx.NewError(code.UserNotFound, `用户不存在`)
		}
		return err
	}
	if len(detail.Mobile) > 0 {
		detail.Mobile = com.MaskString(detail.Mobile)
	}
	if len(detail.Email) > 0 {
		emailInfo := strings.SplitN(detail.Email, `@`, 2)
		emailInfo[0] = com.MaskString(emailInfo[0])
		detail.Email = strings.Join(emailInfo, `@`)
	}
	detail.IdCardNo = ``
	detail.RealName = ``
	detail.SessionId = ``
	if ctx.Format() != echo.ContentTypeHTML {
		mapDetail := detail.AsRow()
		mapDetail.Delete(modelCustomer.PrivateFields...)
		copyDetail := detail.AsMap().Select(`Level`, `Agent`, `Roles`, `Group`).Transform(transformCustomer.Detail)
		mapDetail.DeepMerge(copyDetail)
		ctx.Set(`info`, mapDetail)
	} else {
		ctx.Set(`info`, detail)
	}
	ctx.Internal().Set(`homeowner`, detail)
	blocks := registryUserhome.BlockAll(ctx)
	var block *dashboard.Block
	if len(operate) > 0 {
		for _, v := range blocks {
			if v.Ident == operate {
				block = v
				break
			}
		}
	}
	ctx.Set(`blocks`, blocks)
	if block == nil {
		block = blocks[0]
	}
	if block == nil || block.Hidden.Bool {
		return echo.ErrNotFound
	}
	if err = block.Ready(ctx); err != nil {
		return err
	}
	if ctx.Internal().Bool(`handler.end`) {
		return err
	}
	ctx.Set(`block`, block)
	ctx.Set(`operate`, block.Ident)
	ctx.Set(`operateName`, block.Title)
	ctx.Set(`isUserhome`, true)
	return ctx.Render(`userhome/index`, common.Err(ctx, err))
}

func Homeowner(ctx echo.Context) *modelCustomer.CustomerAndGroup {
	return ctx.Internal().Get(`homeowner`).(*modelCustomer.CustomerAndGroup)
}

func EndHandler(ctx echo.Context) {
	ctx.Internal().Set(`handler.end`, true)
}
