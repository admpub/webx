package official

import (
	"errors"
	"net/url"

	"github.com/admpub/cache"
	"github.com/admpub/license_gen/lib"
	"github.com/admpub/log"
	"github.com/coscms/sms/providers/aliyun"
	"github.com/coscms/sms/providers/twilio"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"

	dbschemaNging "github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/license"
	"github.com/admpub/nging/v4/application/registry/settings"
	"github.com/admpub/nging/v4/application/registry/upload"

	xOAuth "github.com/admpub/webx/application/handler/frontend/oauth"
	xcache "github.com/admpub/webx/application/library/cache"
	cfgIPFilter "github.com/admpub/webx/application/library/ipfilter"
	modelApi "github.com/admpub/webx/application/model/official/api"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"

	modelDBMgr "github.com/nging-plugins/dbmanager/application/model"
)

var configDefaults = map[string]map[string]*dbschemaNging.NgingConfig{
	`base`: {
		`siteLogo`: {
			Key:         `siteLogo`,
			Label:       `站点LOGO(中)`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `image`,
			Sort:        0,
			Disabled:    `N`,
		},
		`siteLogoLarge`: {
			Key:         `siteLogoLarge`,
			Label:       `站点LOGO(大)`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `image`,
			Sort:        0,
			Disabled:    `N`,
		},
		`siteLogoSmall`: {
			Key:         `siteLogoSmall`,
			Label:       `站点LOGO(小)`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `image`,
			Sort:        0,
			Disabled:    `N`,
		},
		`siteFavicon`: {
			Key:         `siteFavicon`,
			Label:       `收藏夹图标`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `image`,
			Sort:        0,
			Disabled:    `N`,
		},
		`siteName`: {
			Key:         `siteName`,
			Label:       `站点名称`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`siteSlogan`: {
			Key:         `siteSlogan`,
			Label:       `站点口号`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`siteURL`: {
			Key:         `siteURL`,
			Label:       `站点网址`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `text`,
			Sort:        1,
			Disabled:    `N`,
		},
		`siteICPFiling`: {
			Key:         `siteICPFiling`,
			Label:       `ICP备案号`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `text`,
			Sort:        2,
			Disabled:    `N`,
		},
		`siteICPURL`: {
			Key:         `siteICPURL`,
			Label:       `ICP备案网址`,
			Description: ``,
			Value:       `https://beian.miit.gov.cn/`,
			Group:       `base`,
			Type:        `text`,
			Sort:        2,
			Disabled:    `N`,
		},
		`siteClose`: {
			Key:         `siteClose`,
			Label:       `关闭站点`,
			Description: ``,
			Value:       `0`, // 0-开放网站;1-关闭网站;2-仅仅会员可以访问;3-仅仅管理员可以访问
			Group:       `base`,
			Type:        `text`,
			Sort:        3,
			Disabled:    `N`,
		},
		`assetsCDN`: {
			Key:         `assetsCDN`,
			Label:       `网站素材CDN`,
			Description: ``,
			Value:       `{"frontend":"","backend":""}`, // frontend-前台素材;backend-后台素材(前台引用时有效)
			Group:       `base`,
			Type:        `json`,
			Sort:        3,
			Disabled:    `N`,
		},
		`showAnnouncement`: {
			Key:         `showAnnouncement`,
			Label:       `显示公告`,
			Description: ``,
			Value:       `0`,
			Group:       `base`,
			Type:        `text`,
			Sort:        4,
			Disabled:    `N`,
		},
		`siteAnnouncement`: {
			Key:         `siteAnnouncement`,
			Label:       `站点公告`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `html`,
			Sort:        5,
			Disabled:    `N`,
		},
		`verifyCodeLifetime`: {
			Key:         `verifyCodeLifetime`,
			Label:       `信件验证码有效期`,
			Description: `短信/邮件验证码的有效期(分钟)`,
			Value:       `5`,
			Group:       `base`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`statCode`: {
			Key:         `statCode`,
			Label:       `统计代码`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `text`,
			Sort:        30,
			Disabled:    `N`,
		},
		`hashidSalt`: {
			Key:         `hashidSalt`,
			Label:       `HashID盐值`,
			Description: ``,
			Value:       `nging`,
			Group:       `base`,
			Type:        `text`,
			Sort:        30,
			Disabled:    `N`,
		},
		`customerRegister`: {
			Key:         `customerRegister`,
			Label:       `前台用户注册`,
			Description: ``,
			Value:       `open`, //close-关闭注册,open-开放注册,invitation-邀请注册
			Group:       `base`,
			Type:        `text`,
			Sort:        30,
			Disabled:    `N`,
		},
		`customerLogin`: {
			Key:         `customerLogin`,
			Label:       `前台用户登录`,
			Description: ``,
			Value:       `open`, //close-关闭登录,open-开放登录
			Group:       `base`,
			Type:        `text`,
			Sort:        30,
			Disabled:    `N`,
		},
		`addIntegral`: {
			Key:         `addIntegral`,
			Label:       `增加消费积分`,
			Description: ``,
			Value:       `{"login":0.0,"register":0.0}`,
			Group:       `base`,
			Type:        `json`,
			Sort:        30,
			Disabled:    `N`,
		},
		`addExperience`: {
			Key:         `addExperience`,
			Label:       `添加经验`,
			Description: ``,
			Value:       `{"login":0.0,"register":0.0}`,
			Group:       `base`,
			Type:        `json`,
			Sort:        30,
			Disabled:    `N`,
		},
		`customerPermCacheTTL`: {
			Key:         `customerPermCacheTTL`,
			Label:       `前台用户权限缓存时长`,
			Description: ``,
			Value:       `-1`, //-1:禁用缓存；-2:强制更新缓存；>=0:缓存时长(秒)
			Group:       `base`,
			Type:        `text`,
			Sort:        30,
			Disabled:    `N`,
		},
		`comment`: {
			Key:         `comment`,
			Label:       `评论`,
			Description: ``,
			Value:       `open`, //close-关闭评论,open-开放评论无需审核,review-开放评论但需要审核
			Group:       `base`,
			Type:        `text`,
			Sort:        30,
			Disabled:    `N`,
		},
		`defaultHandler`: {
			Key:         `defaultHandler`,
			Label:       `首页路由`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `text`,
			Sort:        30,
			Disabled:    `N`,
		},
		`ipFilter`: {
			Key:         `ipFilter`,
			Label:       `IP过滤`,
			Description: ``,
			Value:       `{"On":false,"PassToken":"","BlockByDefault":false,"TrustProxy":true,"AllowedCountries":[],"BlockedCountries":[],"AllowedIPs":[],"BlockedIPs":[]}`,
			Group:       `base`,
			Type:        `json`,
			Sort:        30,
			Disabled:    `N`,
		},
		`recharge`: {
			Key:         `recharge`,
			Label:       `充值`,
			Description: ``,
			Value:       `{"On":false,"MinAmount":100,"DefaultAmount":100}`,
			Group:       `base`,
			Type:        `json`,
			Sort:        30,
			Disabled:    `N`,
		},
		`defaultCurrency`: {
			Key:         `defaultCurrency`,
			Label:       `充值`,
			Description: ``,
			Value:       `CNY`,
			Group:       `base`,
			Type:        `text`,
			Sort:        30,
			Disabled:    `N`,
		},
		`jwtMaxLifetime`: {
			Key:         `jwtMaxLifetime`,
			Label:       `JWT有效期`,
			Description: `JWT有效期(单位:秒)`,
			Value:       `0`,
			Group:       `base`,
			Type:        `text`,
			Sort:        30,
			Disabled:    `N`,
		},
	},
	`cache`: {
		`default`: { // 主要缓存
			Key:         `default`,
			Label:       `主要缓存`,
			Description: ``,
			Value:       `{"Adapter":"redis","AdapterConfig":"","Interval":300,"OccupyMode":false,"Section":""}`,
			Group:       `cache`,
			Type:        `json`,
			Sort:        1000,
			Disabled:    `N`,
		},
		`fallback`: { // 备用缓存
			Key:         `fallback`,
			Label:       `备用缓存`,
			Description: ``,
			Value:       `{"Adapter":"redis","AdapterConfig":"","Interval":300,"OccupyMode":false,"Section":""}`,
			Group:       `cache`,
			Type:        `json`,
			Sort:        1000,
			Disabled:    `N`,
		},
	},
	`contact`: {
		`siteEmail`: {
			Key:         `siteEmail`,
			Label:       `投诉邮箱`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1000,
			Disabled:    `N`,
		},
		`siteMobile`: {
			Key:         `siteMobile`,
			Label:       `投诉电话`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1001,
			Disabled:    `N`,
		},
		`siteQQ`: {
			Key:         `siteQQ`,
			Label:       `投诉QQ`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1002,
			Disabled:    `N`,
		},
		`customerServiceEmail`: {
			Key:         `customerServiceEmail`,
			Label:       `客服邮箱`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1010,
			Disabled:    `N`,
		},
		`customerServiceMobile`: {
			Key:         `customerServiceMobile`,
			Label:       `客服电话`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1011,
			Disabled:    `N`,
		},
		`customerServiceFax`: {
			Key:         `customerServiceFax`,
			Label:       `传真号码`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1012,
			Disabled:    `N`,
		},
		`customerServiceQQ`: {
			Key:         `customerServiceQQ`,
			Label:       `客服QQ`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1013,
			Disabled:    `N`,
		},
		`customerServiceTime`: {
			Key:         `customerServiceTime`,
			Label:       `客服时间`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1014,
			Disabled:    `N`,
		},
		`customerServiceAddress`: {
			Key:         `customerServiceAddress`,
			Label:       `联系地址`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1015,
			Disabled:    `N`,
		},
		`recvMoneyMethod`: {
			Key:         `recvMoneyMethod`,
			Label:       `收款方式(银行名称)`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1050,
			Disabled:    `N`,
		},
		`recvMoneyBranch`: {
			Key:         `recvMoneyBranch`,
			Label:       `收款银行分行(收款方式为银行类时有效)`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1051,
			Disabled:    `N`,
		},
		`recvMoneyAccount`: {
			Key:         `recvMoneyAccount`,
			Label:       `收款账号`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1052,
			Disabled:    `N`,
		},
		`recvMoneyOwner`: {
			Key:         `recvMoneyOwner`,
			Label:       `收款人户名`,
			Description: ``,
			Value:       ``,
			Group:       `contact`,
			Type:        `text`,
			Sort:        1053,
			Disabled:    `N`,
		},
	},

	`oauth`: {
		`alipay`: {
			Key:         `alipay`,
			Label:       `支付宝登录`,
			Description: ``,
			Value:       ``,
			Group:       `oauth`,
			Type:        `json`,
			Sort:        2000,
			Disabled:    `Y`,
		},
		`wechat`: {
			Key:         `wechat`,
			Label:       `微信登录`,
			Description: ``,
			Value:       ``,
			Group:       `oauth`,
			Type:        `json`,
			Sort:        2001,
			Disabled:    `Y`,
		},
		`qq`: {
			Key:         `qq`,
			Label:       `QQ登录`,
			Description: ``,
			Value:       ``,
			Group:       `oauth`,
			Type:        `json`,
			Sort:        2002,
			Disabled:    `Y`,
		},
		`github`: {
			Key:         `github`,
			Label:       `Github登录`,
			Description: ``,
			Value:       ``,
			Group:       `oauth`,
			Type:        `json`,
			Sort:        2003,
			Disabled:    `Y`,
		},
		`microsoft`: {
			Key:         `microsoft`,
			Label:       `Microsoft登录`,
			Description: ``,
			Value:       ``,
			Group:       `oauth`,
			Type:        `json`,
			Sort:        2004,
			Disabled:    `Y`,
		},
	},
	`thirdparty`: { // 第三方接口控制
		`payment`: {
			Key:         `payment`,
			Label:       `收银台远程接口`,
			Description: `支付功能是否使用收银台远程接口（不启用则使用本地配置的付款账号）`,
			Value:       `0`, // 接口账号ID
			Group:       `thirdparty`,
			Type:        `id`,
			Sort:        3000,
			Disabled:    `N`,
		},
		`oauth`: {
			Key:         `oauth`,
			Label:       `社区登录使用远程接口`,
			Description: `社区登录是否使用远程接口（不启用则使用本地配置的社区账号）`,
			Value:       `0`,
			Group:       `thirdparty`,
			Type:        `id`,
			Sort:        3000,
			Disabled:    `N`,
		},
		`exchangeRate`: {
			Key:         `exchangeRate`,
			Label:       `汇率查询接口`,
			Description: `汇率查询接口`,
			Value:       `{"provider":"","apiKey":""}`,
			Group:       `thirdparty`,
			Type:        `json`,
			Sort:        3000,
			Disabled:    `N`,
		},
	},
	`sms`: {
		`aliyun`: {
			Key:         `aliyun`,
			Label:       `阿里云通信`,
			Description: ``,
			Value:       ``,
			Group:       `sms`,
			Type:        `json`,
			Sort:        2000,
			Disabled:    `Y`,
		},
		`twilio`: {
			Key:         `twilio`,
			Label:       `Twilio`,
			Description: ``,
			Value:       ``,
			Group:       `sms`,
			Type:        `json`,
			Sort:        2000,
			Disabled:    `Y`,
		},
	},
	`frequency`: { //频率
		`rateLimiter`: {
			Key:         `rateLimiter`,
			Label:       `限流`,
			Description: ``,
			Value:       `{"On":false,"Max":0,"Duration":0,"Prefix":"","SkipInternalError":true,"RedisAddr":"","RedisPassword":"","RedisDB":0,"DBAccountID":0}`,
			Group:       `frequency`,
			Type:        `text`,
			Sort:        30,
			Disabled:    `N`,
		},
		`mobile`: {
			Key:         `mobile`,
			Label:       `短信频率`,
			Description: `短信的发送频率控制`,
			Value:       `{"maxPerDay":10,"interval":60}`,
			Group:       `frequency`,
			Type:        `json`,
			Sort:        2000,
			Disabled:    `N`,
		},
		`email`: {
			Key:         `email`,
			Label:       `邮件频率`,
			Description: `邮件的发送频率控制`,
			Value:       `{"maxPerDay":10,"interval":60}`,
			Group:       `frequency`,
			Type:        `json`,
			Sort:        2000,
			Disabled:    `N`,
		},
		`message`: {
			Key:         `message`,
			Label:       `发送私信`,
			Description: `私信发送频率控制`,
			Value:       `{"maxPerDay":100,"interval":60}`,
			Group:       `frequency`,
			Type:        `json`,
			Sort:        2000,
			Disabled:    `N`,
		},
		`comment`: {
			Key:         `comment`,
			Label:       `发布评论`,
			Description: `评论发布频率控制`,
			Value:       `{"maxPerDay":100,"maxPending":10}`,
			Group:       `frequency`,
			Type:        `json`,
			Sort:        2000,
			Disabled:    `N`,
		},
		`article`: {
			Key:         `article`,
			Label:       `发布文章`,
			Description: `文章发布频率控制`,
			Value:       `{"maxPerDay":100,"maxPending":10}`,
			Group:       `frequency`,
			Type:        `json`,
			Sort:        2000,
			Disabled:    `N`,
		},
	},
	`custom`: { //独立页面配置
		`login`: {
			Key:         `login`,
			Label:       `登录页面`,
			Description: `登录页面左侧区域内容`,
			Value:       ``,
			Group:       `custom`,
			Type:        `text`,
			Sort:        2000,
			Disabled:    `N`,
		},
	},
	`socketio`: {
		`enabled`: {
			Key:         `enabled`,
			Label:       `socketio开关`,
			Description: `控制socket.io服务是否开启`,
			Value:       ``,
			Group:       `socketio`,
			Type:        `text`,
			Sort:        2000,
			Disabled:    `N`,
		},
	},
}

func init() {
	// 添加默认配置数据
	settings.AddConfigs(configDefaults)

	// 注册配置模板和逻辑
	settings.AddTmpl(`base`, `official/settings/base`, settings.OptAddHookGet(func(ctx echo.Context) error {
		var storerEngines []string
		for name := range upload.StorerAll() {
			storerEngines = append(storerEngines, name)
		}
		ctx.Set(`storerEngines`, storerEngines)
		return nil
	}))
	settings.Register((&settings.SettingForm{
		Short: `缓存配置`,
		Label: `缓存配置`,
		Group: `cache`,
		Tmpl:  []string{`official/settings/cache`},
	}).AddHookGet(func(ctx echo.Context) error {
		dbaM := modelDBMgr.NewDbAccount(ctx)
		dbaM.ListByOffset(nil, nil, 0, -1, db.Cond{`engine`: `redis`})
		ctx.Set(`dbAccounts`, dbaM.Objects())
		ctx.Set(`cacheAdapters`, []string{`redis`, `file`})
		ctx.SetFunc(`isDbAccount`, xcache.IsDbAccount)
		return nil
	}))
	settings.Register((&settings.SettingForm{
		Short: `短信服务`,
		Label: `短信服务`,
		Group: `sms`,
		Tmpl:  []string{`official/settings/sms`},
	}).AddHookPost(func(_ echo.Context) error {
		return xOAuth.UpdateSMSConfigs()
	}))
	settings.Register((&settings.SettingForm{
		Short: `频率控制`,
		Label: `频率控制`,
		Group: `frequency`,
		Tmpl:  []string{`official/settings/frequency`},
	}).AddHookGet(func(ctx echo.Context) error {
		dbaM := modelDBMgr.NewDbAccount(ctx)
		dbaM.ListByOffset(nil, nil, 0, -1, db.Cond{`engine`: `redis`})
		ctx.Set(`dbAccounts`, dbaM.Objects())
		return nil
	}))
	settings.Register((&settings.SettingForm{
		Short: `联系方式`,
		Label: `官方联系方式`,
		Group: `contact`,
		Tmpl:  []string{`official/settings/contact`},
	}))
	settings.Register((&settings.SettingForm{
		Short: `第三方接口`,
		Label: `第三方接口设置`,
		Group: `thirdparty`,
		Tmpl:  []string{`official/settings/thirdparty`},
	}).AddHookGet(func(ctx echo.Context) error {
		m := modelApi.NewAccount(ctx)
		m.ListByOffset(nil, nil, 0, -1, db.Cond{`disabled`: `N`})
		ctx.Set(`apiAccounts`, m.Objects())
		return nil
	}))
	settings.Register((&settings.SettingForm{
		Short: `第三方登录`,
		Label: `第三方登录接口设置`,
		Group: `oauth`,
		Tmpl:  []string{`official/settings/oauth`},
	}).AddHookPost(func(_ echo.Context) error {
		return xOAuth.UpdateAccount()
	}))
	settings.Register((&settings.SettingForm{
		Short: `自定义页面`,
		Label: `自定义页面设置`,
		Group: `custom`,
		Tmpl:  []string{`official/settings/custom`},
	}))
	settings.Register((&settings.SettingForm{
		Short: `SocketIO`,
		Label: `SocketIO设置`,
		Group: `socketio`,
		Tmpl:  []string{`official/settings/socketio`},
	}))

	// 注册配置值解码器（用于从数据库读出来之后的解码操作）
	// 名称支持"group"或"group.key"两种格式(优先使用指定了"group.key"的解码器，否则使用指定了"group"的解码器)，例如:
	// settings.RegisterDecoder(`sms`,...)作为sms组的默认解码器
	// settings.RegisterDecoder(`sms.twilio`,...)对sms组内key为twilio的配置有效
	settings.RegisterDecoder(`base.ipFilter`, func(v *dbschemaNging.NgingConfig, r echo.H) error {
		jsonData := cfgIPFilter.NewOptions()
		if len(v.Value) > 0 {
			com.JSONDecode(com.Str2bytes(v.Value), jsonData)
		}
		r[`ValueObject`] = jsonData
		return nil
	})
	settings.RegisterDecoder(`base.recharge`, func(v *dbschemaNging.NgingConfig, r echo.H) error {
		jsonData := modelCustomer.NewWalletSettings()
		if len(v.Value) > 0 {
			com.JSONDecode(com.Str2bytes(v.Value), jsonData)
		}
		r[`ValueObject`] = jsonData
		return nil
	})
	settings.RegisterDecoder(`frequency.rateLimiter`, func(v *dbschemaNging.NgingConfig, r echo.H) error {
		jsonData := cfgIPFilter.NewRateLimiterConfig()
		if len(v.Value) > 0 {
			com.JSONDecode(com.Str2bytes(v.Value), jsonData)
		}
		r[`ValueObject`] = jsonData
		return nil
	})
	settings.RegisterDecoder(`cache`, func(v *dbschemaNging.NgingConfig, r echo.H) error {
		jsonData := &cache.Options{}
		if len(v.Value) > 0 {
			com.JSONDecode(com.Str2bytes(v.Value), jsonData)
		}
		r[`ValueObject`] = jsonData
		return nil
	})
	settings.RegisterDecoder(`sms`, func(v *dbschemaNging.NgingConfig, r echo.H) error {
		var jsonData interface{}
		switch v.Key {
		case `twilio`:
			jsonData = twilio.New()
		case `aliyun`:
			jsonData = aliyun.New()
		default:
			return errors.New(`The decoder is unsupported: ` + v.Key)
		}
		if len(v.Value) > 0 {
			com.JSONDecode(com.Str2bytes(v.Value), jsonData)
		}
		r[`ValueObject`] = jsonData
		return nil
	})
	license.OnSetLicense(func(data *lib.LicenseData) {
		if !config.IsInstalled() {
			return
		}
		if license.LicenseMode() != license.ModeDomain {
			return
		}
		siteURL := config.Setting(`base`).String(`siteURL`)
		if len(siteURL) == 0 {
			return
		}
		u, err := url.Parse(siteURL)
		if err != nil {
			return
		}
		if license.EqDomain(u.Hostname(), data.Info.Domain) {
			return
		}
		frontendURL := u.Scheme + `://` + data.Info.Domain
		config.Setting(`base`).Set(`siteURL`, frontendURL)
		subdomains.SetBaseURL(`frontend`, frontendURL)
		if config.FromFile() == nil || !config.FromFile().ConnectedDB() {
			return
		}
		v := dbschemaNging.NewNgingConfig(common.NewMockContext())
		err = v.UpdateField(nil, `value`, frontendURL, db.And(
			db.Cond{`key`: `siteURL`},
			db.Cond{`group`: `base`},
		))
		if err != nil {
			panic(err)
		}
	})
	settings.RegisterDecoder(`base.siteURL`, func(v *dbschemaNging.NgingConfig, _ echo.H) error {
		// 从数据库里读出来的时候
		frontendURL := v.Value
		err := license.CheckSiteURL(frontendURL, true)
		if err == nil {
			return nil
		}
		log.Error(err)
		domain := license.FullDomain()
		if len(domain) > 0 {
			u, err := url.Parse(frontendURL)
			if err != nil {
				log.Errorf(`%s: %v`, frontendURL, err.Error())
				return nil
			}
			frontendURL = u.Scheme + `://` + domain
			v.Value = frontendURL
			log.Warn(`reset siteURL: `, v.Value)
			err = v.UpdateField(nil, `value`, v.Value, db.And(
				db.Cond{`key`: v.Key},
				db.Cond{`group`: v.Group},
			))
			if err != nil {
				panic(err)
			}
		}
		return nil
	})

	// 注册配置值编码器（用于客户端提交表单数据之后的编码操作，编码结果保存到数据库）
	settings.RegisterEncoder(`base.ipFilter`, func(v *dbschemaNging.NgingConfig, r echo.H) ([]byte, error) {
		cfg := cfgIPFilter.NewOptions().FromStore(r)
		return com.JSONEncode(cfg)
	})
	settings.RegisterEncoder(`base.recharge`, func(v *dbschemaNging.NgingConfig, r echo.H) ([]byte, error) {
		cfg := modelCustomer.NewWalletSettings().FromStore(r)
		return com.JSONEncode(cfg)
	})
	settings.RegisterEncoder(`frequency.rateLimiter`, func(v *dbschemaNging.NgingConfig, r echo.H) ([]byte, error) {
		cfg := cfgIPFilter.NewRateLimiterConfig().FromStore(r)
		return com.JSONEncode(cfg)
	})
	settings.RegisterEncoder(`cache`, func(v *dbschemaNging.NgingConfig, r echo.H) ([]byte, error) {
		cfg := &cache.Options{
			Adapter:       r.String(`Adapter`),
			AdapterConfig: r.String(`AdapterConfig`),
			Interval:      r.Int(`Interval`),
			OccupyMode:    r.Bool(`OccupyMode`),
		}
		return com.JSONEncode(cfg)
	})
	settings.RegisterEncoder(`sms`, func(v *dbschemaNging.NgingConfig, r echo.H) ([]byte, error) {
		switch v.Key {
		case `twilio`:
			cfg := twilio.New()
			cfg.AccessKey = r.String(`accessKey`)
			cfg.AccessSecret = r.String(`accessSecret`)
			cfg.From = r.String(`from`)
			cfg.CountryCode = r.String(`countryCode`)
			return com.JSONEncode(cfg)
		case `aliyun`:
			cfg := aliyun.New()
			cfg.GatewayURL = r.String(`gatewayURL`, cfg.GatewayURL)
			cfg.AccessKey = r.String(`accessKey`)
			cfg.AccessSecret = r.String(`accessSecret`)
			cfg.SignName = r.String(`signName`) //默认签名
			cfg.TmplCode = r.String(`tmplCode`) //默认模板代码
			return com.JSONEncode(cfg)
		default:
			return nil, errors.New(`The decoder is unsupported: ` + v.Key)
		}
	})
	settings.RegisterEncoder(`base.siteURL`, func(v *dbschemaNging.NgingConfig, formDataMap echo.H) ([]byte, error) {
		//写入数据库的时候
		frontendURL := formDataMap.String(`value`)
		err := license.CheckSiteURL(frontendURL, true)
		return []byte(frontendURL), err
	})

	var updateFrontendURL = func(diff config.Diff) (err error) {
		//设置内存变量的时候
		frontendURL := diff.String()
		subdomains.SetBaseURL(`frontend`, frontendURL)
		return
	}
	config.OnKeySetSettings(`base.siteURL`, updateFrontendURL)
}
