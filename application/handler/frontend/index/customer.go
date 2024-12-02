package index

import (
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/session"

	"github.com/admpub/log"
	"github.com/admpub/webx/application/handler/frontend/user"
	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/ip2region"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webcore/library/sessionguard"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/top"
	"github.com/coscms/webfront/middleware/sessdata"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

// CookieMaxAge 允许设置的Cookie最大有效时长(单位:秒)
var CookieMaxAge = 86400 * 365

const debug = false

func SetJWTData(c echo.Context, m *modelCustomer.Customer) error {
	enabledJWT := c.Internal().Bool(`enabledJWT`)
	if !enabledJWT {
		return nil
	}
	data := c.Data()
	signed, err := m.JWTSignedString(nil)
	if err != nil {
		data.SetError(err)
	} else {
		newData := echo.H{
			`jwt`: signed,
			`sid`: c.Session().MustID(),
		}
		if m.OfficialCustomer.Id > 0 {
			newData[`profile`] = m.ClearPasswordData()
		}
		h, y := data.GetData().(echo.H)
		if !y {
			data.SetData(newData)
		} else {
			h.DeepMerge(newData)
		}
	}
	return nil
}

// SignUp 注册
func SignUp(c echo.Context) error {
	var err error

	next := c.Form(`next`)
	next = echo.GetOtherURL(c, next)
	if len(next) == 0 {
		next = sessdata.URLFor(`/index`)
	}

	// 已经登录的时候跳过当前页面
	if sessdata.Customer(c) != nil {
		c.Data().SetInfo(c.T(`您已经是登录状态了, 无需注册`))
		return c.Redirect(next)
	}
	registerState := config.Setting(`base`).String(`customerRegister`)
	if c.IsPost() {
		m := modelCustomer.NewCustomer(c)
		email := c.Form(`email`)
		mobile := c.Form(`mobile`)
		name := c.Form(`name`)
		pass := c.Form(`password`)
		var onReg func(m *modelCustomer.Customer) error
		switch registerState {
		case `close`:
			err = c.E(`用户注册功能暂时关闭`)
			goto END
		case `invitation`:
			invitationCode := c.Form(`invitation`)
			invM := modelCustomer.NewInvitation(c)
			err = invM.FindCode(invitationCode)
			if err != nil {
				goto END
			}
			m.LevelId = invM.LevelId
			m.AgentLevel = invM.AgentLevelId
			m.RoleIds = invM.RoleIds
			onReg = func(m *modelCustomer.Customer) error {
				return invM.UseCode(invM.Id, m.OfficialCustomer)
			}
		default: //open
		}
		if pass != c.Form(`repassword`) {
			err = c.NewError(code.InvalidParameter, `两次输入的密码不匹配`).SetZone(`repassword`)
		} else if data := captchabiz.VerifyCaptcha(c, httpserver.KindFrontend, `code`); nerrors.IsFailureCode(data.GetCode()) {
			if err := SetJWTData(c, m); err != nil {
				return err
			}
			err = c.NewError(code.InvalidParameter, `验证码不正确`).SetZone(`code`)
			if c.Format() == `json` {
				return c.JSON(data)
			}
		} else {
			if err = c.Begin(); err != nil {
				goto END
			}
			m.SetContext(c)
			err = m.SignUp(name, pass, mobile, email, modelCustomer.GenerateOptionsFromHeader(c)...)
			if err != nil {
				c.Rollback()
				goto END
			}
			if onReg != nil {
				if err = onReg(m); err != nil {
					c.Rollback()
					goto END
				}
			}
			c.Commit()
			if err := SetJWTData(c, m); err != nil {
				return err
			}
			return c.Redirect(next)
		}
	}

END:
	tmpl := c.Internal().String(`tmpl`)
	if len(tmpl) == 0 {
		tmpl = `sign_up`
	}
	signInURL := c.Internal().String(`signInURL`)
	if len(signInURL) == 0 {
		signInURL = sessdata.URLFor(`/sign_in`)
	}
	c.Set(`signInURL`, signInURL)
	c.Set(`registerState`, registerState)
	return c.Render(tmpl, common.Err(c, err))
}

// SignOut 退出登录
func SignOut(c echo.Context) error {
	m := modelCustomer.NewCustomer(c)
	deviceM := modelCustomer.NewDevice(c)
	var err error
	var copied dbschema.OfficialCustomer
	customer := sessdata.Customer(c)
	if customer == nil {
		goto END
	}
	copied = *customer
	m.OfficialCustomer = &copied
	m.SetContext(c)
	deviceM.SessionId = c.Session().ID()
	deviceM.CustomerId = customer.Id
	err = deviceM.CleanCustomer(customer, modelCustomer.GenerateOptionsFromHeader(c)...)
	if err != nil {
		log.Error(err)
	}
	m.UnsetSession()

END:
	return c.Redirect(sessdata.URLFor(`/sign_in`))
}

// SignIn 登录
func SignIn(c echo.Context) error {
	var err error
	if c.Formx(`modal`).Bool() {
		tmpl := c.Internal().String(`modalTmpl`)
		if len(tmpl) == 0 {
			tmpl = `partial_modal_sign_in`
		}
		return c.Render(tmpl, err)
	}

	next := c.Form(`next`)
	next = echo.GetOtherURL(c, next)
	if len(next) == 0 {
		next = sessdata.URLFor(`/index`)
	}
	// 已经登录的时候跳过当前页面
	if !debug && sessdata.Customer(c) != nil {
		c.Data().SetInfo(c.T(`已经登录过了`))
		return c.Redirect(next)
	}

	if c.IsPost() {
		m := modelCustomer.NewCustomer(c)
		name := c.Form(`name`)
		pass := c.Form(`password`)
		typi := c.Form(`type`)

		data := captchabiz.VerifyCaptcha(c, httpserver.KindFrontend, `code`)
		if !debug && nerrors.IsFailureCode(data.GetCode()) {
			if err := SetJWTData(c, m); err != nil {
				return err
			}
			err = c.NewError(code.InvalidParameter, `验证码不正确`).SetZone(`code`)
			if c.Format() == `json` {
				return c.JSON(data)
			}
		} else {
			remember := c.Form(`remember`)
			var maxAge int
			if len(remember) > 0 {
				if remember == `forever` {
					maxAge = CookieMaxAge
				} else {
					duration, err := top.ParseDuration(remember)
					if err == nil {
						maxAge = int(duration.Seconds())
					}
				}
				if maxAge > CookieMaxAge {
					maxAge = CookieMaxAge
				}
				session.RememberMaxAge(c, maxAge)
			}
			err = m.SignIn(name, pass, typi, modelCustomer.GenerateOptionsFromHeader(c, maxAge)...)
			if err != nil {
				if c.Format() != `html` {
					return err
				}
				goto END
			}
			data.SetInfo(c.T(`登录成功`))
			if err := SetJWTData(c, m); err != nil {
				return err
			}
			return c.Redirect(next)
		}
	}

END:
	tmpl := c.Internal().String(`tmpl`)
	if len(tmpl) == 0 {
		tmpl = `sign_in`
	}
	signUpURL := c.Internal().String(`signUpURL`)
	if len(signUpURL) == 0 {
		signUpURL = sessdata.URLFor(`/sign_up`)
	}
	c.Set(`signUpURL`, signUpURL)
	forgotURL := c.Internal().String(`forgotURL`)
	if len(forgotURL) == 0 {
		forgotURL = sessdata.URLFor(`/forgot`)
	}
	c.Set(`forgotURL`, forgotURL)
	return c.Render(tmpl, common.Err(c, err))
}

func CustomerInfo(c echo.Context) error {
	customer := sessdata.Customer(c)
	data := c.Data()
	if customer == nil {
		data.SetError(c.NewError(code.Unauthenticated, ``))
	} else {
		data.SetData(customer.AsRow(`id`, `uid`, `group_id`, `name`, `online`, `gender`, `avatar`, `role_ids`))
	}
	return c.JSON(data)
}

func qrcodeSignIn(ctx echo.Context) error {
	expireTime := time.Now().Add(time.Minute * 10)
	signInData := &user.QRSignIn{
		SessionID: ctx.Session().MustID(),
		Expires:   expireTime.Unix(),
		Environment: sessionguard.Environment{
			UserAgent: ctx.Request().UserAgent(),
		},
	}
	signInData.Environment.Location, _ = ip2region.IPInfo(ctx.RealIP())
	plaintext, err := com.JSONEncodeToString(signInData)
	if err != nil {
		return err
	}
	qrcode := config.FromFile().Encode256(plaintext)
	qrcode = com.URLSafeBase64(qrcode, true)
	return ctx.JSON(ctx.Data().SetData(echo.H{
		`qrcode`:  qrcode,
		`expires`: expireTime.Format(time.DateTime),
	}))
}
