package oauth

import (
	"strings"

	"github.com/admpub/goth"
	"github.com/coscms/webcore/library/config"
	"github.com/admpub/webx/application/library/cache"
	"github.com/admpub/webx/application/library/oauth2client/providers/wechatgh"
	"github.com/admpub/webx/application/library/thirdparty/login/wechat/gh"
	"github.com/admpub/webx/application/middleware/sessdata"
	oconfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/silenceper/wechat/v2/officialaccount/server"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func init() {
	gh.RegisterEventHandler(gh.EventSubscribe, wechatGHSubscribe)
	gh.RegisterEventHandler(gh.EventScan, wechatGHSubscribe)
}

func makeSceneStr(ctx echo.Context, nonce string) string {
	return `oauth@` + com.Md5(ctx.Session().MustID())[0:16] + nonce
}

// WechatGH 微信公众号登录
func WechatGH(ctx echo.Context) error {
	// 已经登录的时候跳过当前页面
	if sessdata.Customer(ctx) != nil {
		ctx.Data().SetInfo(ctx.T(`已经登录过了`))
		next := ctx.Form(`next`)
		next = echo.GetOtherURL(ctx, next)
		if len(next) == 0 {
			next = sessdata.URLFor(`/index`)
		}
		return ctx.Redirect(next)
	}
	cfg := config.Setting(`oauth`, `wechat`)
	if len(cfg) == 0 {
		return ctx.NewError(code.Unsupported, ctx.T(`不支持此种登录方式`))
	}
	if ctx.IsPost() {
		data := ctx.Data()
		nonce := ctx.Form(`nonce`)
		if len(nonce) != 8 {
			return ctx.NewError(code.InvalidParameter, ``).SetZone(`nonce`)
		}
		sceneStr := makeSceneStr(ctx, nonce)
		var openID string
		err := cache.Get(ctx, sceneStr, &openID)
		if err != nil {
			if cache.IsNotExist(err) {
				err = ctx.NewError(code.DataNotFound, ``)
			}
			data.SetError(err)
			return ctx.JSON(data)
		}
		ocfg, err := newWechatGHConfig(ctx, cfg)
		if err != nil {
			return err
		}
		if err := cache.Delete(ctx, sceneStr); err != nil {
			ctx.Logger().Warn(err)
		}
		wechatUser, err := wechatgh.GetAccount(ocfg).GetUser().GetUserInfo(openID)
		if err != nil {
			return err
		}
		gothUser := goth.User{
			UserID:    wechatUser.OpenID,
			NickName:  wechatUser.Nickname,
			Name:      wechatUser.Nickname,
			AvatarURL: wechatUser.Headimgurl,
			Provider:  `wechatgh`,
			RawData: map[string]interface{}{
				`gender`:  wechatUser.Sex,
				`unionid`: wechatUser.UnionID,
			},
			Location: wechatUser.City + `,` + wechatUser.Province + `,` + wechatUser.Country,
		}
		ctx.Internal().Set(Default().Config.ContextKey, gothUser)
		return successHandler(ctx)
	}
	ocfg, err := newWechatGHConfig(ctx, cfg)
	if err != nil {
		return err
	}
	svr := wechatgh.GetServer(ctx, ocfg)
	accessToken, err := svr.GetAccessToken()
	if err != nil {
		return err
	}
	nonce := com.RandomAlphanumeric(8)
	sceneStr := makeSceneStr(ctx, nonce)
	qrcodeURL, err := gh.GetQRCodeURL(accessToken, sceneStr, 3600)
	if err != nil {
		return err
	}
	data := ctx.Data()
	result := echo.H{
		`qrcodeURL`: qrcodeURL,
		`nonce`:     nonce,
	}
	return ctx.JSON(data.SetData(result))
}

// WechatGHCheckSign 微信公众号签名检查
func WechatGHCheckSign(ctx echo.Context) error {
	data := ctx.Data()
	cfg := config.Setting(`oauth`, `wechat`)
	if len(cfg) == 0 {
		err := ctx.NewError(code.Unsupported, `不支持此种登录方式`)
		return ctx.JSON(data.SetError(err))
	}
	token := cfg.String(`token`)
	if len(token) == 0 {
		err := ctx.NewError(code.DataUnavailable, `token 未配置`)
		return ctx.JSON(data.SetError(err))
	}
	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	echostr := ctx.Query("echostr")
	if len(signature) == 0 || len(timestamp) == 0 || len(nonce) == 0 || len(echostr) == 0 {
		err := ctx.NewError(code.InvalidParameter, ``)
		return ctx.JSON(data.SetError(err))
	}
	ok := wechatgh.CheckSignature(signature, timestamp, nonce, token)
	if !ok {
		return echo.ErrBadRequest
	}

	return ctx.String(echostr)
}

// WechatGHCallback 微信公众号回调
func WechatGHCallback(ctx echo.Context) error {
	data := ctx.Data()
	cfg := config.Setting(`oauth`, `wechat`)
	if len(cfg) == 0 {
		err := ctx.NewError(code.Unsupported, `不支持此种登录方式`)
		return ctx.JSON(data.SetError(err))
	}
	ocfg, err := newWechatGHConfig(ctx, cfg)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	return wechatgh.MessageSystem(ctx, ocfg, gh.MakeMessageHandler(ctx, nil))
}

func newWechatGHConfig(ctx echo.Context, cfg echo.H) (*oconfig.Config, error) {
	token := cfg.String(`token`)
	if len(token) == 0 {
		return nil, ctx.NewError(code.DataUnavailable, `token 未配置`)
	}
	ocfg := &oconfig.Config{
		AppID:          cfg.String(`key`),
		AppSecret:      cfg.String(`secret`),
		Token:          token,
		EncodingAESKey: cfg.String(`encodingAESKey`),
	}
	return ocfg, nil
}

func wechatGHSubscribe(c echo.Context, _ *server.Server, m *message.MixMessage) *message.Reply {
	sceneStr := m.EventKey
	if len(sceneStr) == 0 {
		return nil
	}
	tmpSlice := strings.SplitN(sceneStr, `_`, 2)
	if len(tmpSlice) == 2 {
		sceneStr = tmpSlice[1]
	}
	if strings.HasPrefix(sceneStr, `oauth@`) {
		err := cache.Put(c, sceneStr, m.FromUserName, 3600)
		if err != nil {
			c.Logger().Error(err)
		}
	}
	return nil
}
