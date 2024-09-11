package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/license"
	uploadLibrary "github.com/coscms/webcore/library/upload"
	nav "github.com/coscms/webcore/registry/navigate"

	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/handler/frontend/user/navigate"
	"github.com/admpub/webx/application/library/cache"
	"github.com/admpub/webx/application/library/oauth2client"
	"github.com/admpub/webx/application/library/xconst"
	"github.com/admpub/webx/application/library/xrole"
	"github.com/admpub/webx/application/library/xrole/xroleutils"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func SessionInfo(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ppath := c.Request().URL().Path()
		switch ppath {
		case c.Echo().Prefix() + `/favicon.ico`:
			return h.Handle(c)
		default:
			if strings.HasPrefix(ppath, c.Echo().Prefix()+uploadLibrary.UploadURLPath) {
				return h.Handle(c)
			}
		}
		baseCfg := config.Setting(`base`)
		siteClose := baseCfg.Uint(`siteClose`)
		if siteClose == xconst.SiteClose {
			return close(c, baseCfg)
		}
		var customer *dbschema.OfficialCustomer
		enabledJWT := c.Internal().Bool(`enabledJWT`)
		if enabledJWT {
			customerM := modelCustomer.NewCustomer(c)
			var err error
			customer, err = customerM.GetByJWT()
			if err != nil {
				log.Debug(err.Error())
			}
		}
		if customer == nil {
			customer, _ = c.Session().Get(`customer`).(*dbschema.OfficialCustomer)
		}
		if customer != nil {
			if siteClose == xconst.SiteOnlyAdmin && customer.Uid < 1 {
				return close(c, baseCfg)
			}
			c.Internal().Set(`customer`, customer)
		} else {
			if siteClose == xconst.SiteOnlyMember || siteClose == xconst.SiteOnlyAdmin {
				switch path.Base(ppath) {
				case `sign_up`, `sign_in`, `sign_out`:
				default:
					if c.Path() != c.Echo().Prefix()+`/captcha/*` && !strings.HasPrefix(ppath, c.Echo().Prefix()+`/oauth/`) {
						return goToSignIn(c)
					}
				}
			}
			ouser, exists, err := oauth2client.GetSession(c)
			//echo.Dump(echo.H{`err`: err, `ouser`: ouser})
			if err == nil {
				c.Set(`oauth`, ouser) // 表单数据
			} else {
				if exists {
					log.Debug(err.Error())
				}
			}
		}
		return h.Handle(c)
	}
}

func close(c echo.Context, baseCfg echo.H) error {
	showAnnouncement := baseCfg.Bool(`showAnnouncement`)
	if !showAnnouncement {
		return c.Render(`error/under-maintenance`, nil, http.StatusServiceUnavailable)
		//return c.Render(`error/not-found`, nil, http.StatusNotFound)
	}
	siteAnnouncement := baseCfg.String(`siteAnnouncement`)
	siteAnnouncement = strings.TrimSpace(siteAnnouncement)
	data := echo.H{
		`title`:         `Coming Soon`,
		`content`:       siteAnnouncement,
		`enabledNotify`: false, //是否支持访客接收邮件通知(未实现)
	}
	if strings.HasPrefix(siteAnnouncement, `<h1>`) {
		siteAnnouncement = strings.TrimPrefix(siteAnnouncement, `<h1>`)
		splited := strings.SplitN(siteAnnouncement, `</h1>`, 2)
		switch len(splited) {
		case 2:
			data[`title`] = splited[0]
			data[`content`] = splited[1]
		}
	}
	return c.Render(`error/coming-soon`, data)
}

func userCenter(c echo.Context, customer *dbschema.OfficialCustomer) error {
	m := modelCustomer.NewCustomer(c)
	err := m.VerifySession(customer)
	if err != nil {
		if common.IsUserNotLoggedIn(err) {
			common.SendErr(c, err)
			return goToSignIn(c)
		}
		return err
	}
	return permCheck(c, customer)
}

func permCheck(c echo.Context, customer *dbschema.OfficialCustomer) error {
	permission := xroleutils.CustomerPermission(c, customer)
	//echo.Dump(permission)
	customerID := fmt.Sprint(customer.Id)
	cacheTTL := xroleutils.CustomerPermTTL(c)
	ttlOpt := cache.GetTTLByNumber(cacheTTL, nil)
	c.SetFunc(`LeftNavigate`, func() nav.List {
		list := &nav.List{}
		cache.XFunc(c, sessdata.LeftNavigateCacheKey+customerID, list, func() error {
			*list = permission.FilterNavigate(c, navigate.LeftNavigate)
			return nil
		}, ttlOpt)
		return *list
	})
	c.SetFunc(`TopNavigate`, func() nav.List {
		list := &nav.List{}
		cache.XFunc(c, sessdata.TopNavigateCacheKey+customerID, list, func() error {
			*list = permission.FilterNavigate(c, navigate.TopNavigate)
			return nil
		}, ttlOpt)
		return *list
	})
	if !c.Internal().Bool(`skipCurrentURLPermCheck`) {
		rpath := c.Path()
		if len(c.Echo().Prefix()) > 0 {
			rpath = strings.TrimPrefix(rpath, c.Echo().Prefix())
		}
		if err := checkPermission(c, customer, permission, rpath); err != nil {
			return err
		}
	}
	c.SetFunc(`CheckPerm`, func(route string) error {
		return checkPermission(c, customer, permission, route)
	})
	return nil
}

func checkPermission(ctx echo.Context, customer *dbschema.OfficialCustomer, permission *xrole.RolePermission, route string) error {
	var (
		err error
		ret bool
	)
	err, route, ret = xrole.SpecialAuths.Check(ctx, route, customer, permission)
	if ret || err != nil {
		return err
	}
	if route == `/user/index` {
		return nil
	}
	route = strings.TrimPrefix(route, `/user/`)
	if !permission.Check(ctx, route) {
		return common.ErrUserNoPerm
	}
	return nil
}

func SkipCurrentURLPermCheck(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Internal().Set(`skipCurrentURLPermCheck`, true)
		return h.Handle(c)
	}
}

func AuthCheck(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		//检查是否已安装
		if !config.IsInstalled() {
			c.Data().SetError(c.NewError(stdCode.SystemNotInstalled, c.T(`请先安装`)))
			return c.Redirect(subdomains.Default.URL(`/setup`, `backend`))
		}

		//验证授权文件
		if !license.Ok(c) {
			c.Data().SetError(c.NewError(stdCode.SystemUnauthorized, c.T(`请先获取本系统授权`)))
			return c.Redirect(subdomains.Default.URL(`/license`, `backend`))
		}

		customer := Customer(c)
		if customer != nil {
			if err := userCenter(c, customer); err != nil {
				return err
			}
			return h.Handle(c)
		}
		return goToSignIn(c)
	}
}

func TrimPathSuffix(ignorePrefixes ...string) echo.MiddlewareFuncd {
	return func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			upath := c.Request().URL().Path()
			for _, ignorePrefix := range ignorePrefixes {
				if strings.HasPrefix(upath, ignorePrefix) {
					return h.Handle(c)
				}
			}
			cleanedPath := strings.TrimSuffix(upath, c.DefaultExtension())
			c.Request().URL().SetPath(cleanedPath)
			return h.Handle(c)
		}
	}
}

func goToSignIn(c echo.Context) error {
	var queryString string
	if c.IsGet() {
		next := c.Request().URI()
		if !strings.Contains(next, `/sign_in`) {
			queryString = `?next=` + url.QueryEscape(next)
		}
	}
	c.Data().SetError(c.NewError(stdCode.Unauthenticated, c.T(`请先登录`)))
	return c.Redirect(URLFor(`/sign_in`) + queryString)
}

var (
	Customer   = sessdata.Customer
	AgentLevel = sessdata.AgentLevel
	URLFor     = sessdata.URLFor
)

var Middlewares []interface{}

func Use(m ...interface{}) {
	Middlewares = append(Middlewares, m...)
}
