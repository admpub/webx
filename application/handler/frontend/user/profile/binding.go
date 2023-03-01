package profile

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/webx/application/handler/frontend/user/binding"
	xMW "github.com/admpub/webx/application/middleware"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func Binding(ctx echo.Context) (err error) {
	customer := xMW.Customer(ctx)
	typ := ctx.Form(`type`, `email`)
	m := modelCustomer.NewCustomer(ctx)
	err = m.Get(nil, `id`, customer.Id)
	if err != nil {
		return
	}
	binder := binding.Get(typ)
	if binder == nil {
		return ctx.E(`不支持的类型: %s`, typ)
	}
	if ctx.IsPost() {
		verify := ctx.Formx(`verify`).Bool()
		if verify {
			return binder.Verifier(ctx, m)
		}
		return binder.Sender(ctx, m)
	}
	var (
		remainCount    int64
		waitingSeconds int64
		maxPerDay      int64
	)
	if typ == `oauth` {
		if len(ctx.Form(`provider`)) > 0 {
			//跳转到第三方平台进行登录验证
			return binder.Sender(ctx, m)
		}
		//oauth不需要本地验证，所以对于他的定义比较特殊
		//这里用作查询支持绑定的第三方平台账号和当前账号是否已与其绑定
		if err = binder.Verifier(ctx, m); err != nil {
			return err
		}
	} else {
		waitingSeconds = BindingSendingInterval(typ)
		maxPerDay = BindingSendingMaxPerDay(typ)
		remainCount = BindingSendingRemainCount(ctx, m, typ)
	}
	ret := handler.Err(ctx, err)
	ctx.Set(`activeURL`, `/user/profile`)
	ctx.Request().Form().Set(`type`, typ)
	ctx.Set(`objectName`, binder.ObjectName)
	ctx.Set(`typeName`, binder.Name)
	ctx.Set(`waitingSeconds`, waitingSeconds)
	ctx.Set(`maxPerDay`, maxPerDay)
	ctx.Set(`remainCount`, remainCount)
	ctx.Set(`title`, ctx.T(`账号绑定`))
	return ctx.Render(`user/profile/binding`, ret)
}

// BindingSendingInterval 间隔时间(秒)
func BindingSendingInterval(sendMethod string) int64 {
	frequencyCfg := config.Setting().GetStore(`frequency`).GetStore(sendMethod)
	interval := frequencyCfg.Int64(`interval`, model.SMSWaitingSeconds)
	return interval
}

// BindingSendingMaxPerDay 每人每天发送上限
func BindingSendingMaxPerDay(sendMethod string) int64 {
	frequencyCfg := config.Setting().GetStore(`frequency`).GetStore(sendMethod)
	maxPerDay := frequencyCfg.Int64(`maxPerDay`, model.SMSMaxPerDay)
	return maxPerDay
}

// BindingSendingRemainCount 今天剩余可发几次
func BindingSendingRemainCount(ctx echo.Context, cust *modelCustomer.Customer, sendMethod string) int64 {
	maxPerDay := BindingSendingMaxPerDay(sendMethod)
	m := model.NewCode(ctx)
	n, _ := m.CountTodayVerificationCode(cust.Id, `customer`, sendMethod)
	remain := maxPerDay - n
	if remain < 0 {
		return 0
	}
	return remain
}
