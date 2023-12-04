package oauth2client

import (
	"encoding/json"
	"errors"

	"github.com/admpub/goth"
	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

var ErrGetSession = errors.New("get user data from session failed")
var OAuthUserSessionKey = `oauthUser`

func SaveSession(ctx echo.Context, ouser *goth.User) error {
	b, err := json.Marshal(ouser)
	if err != nil {
		return err
	}
	log.Debugf(`[oauth]SaveSession:%s=%s`, OAuthUserSessionKey, string(b))
	ctx.Session().Set(OAuthUserSessionKey, b)
	return nil
}

func GetSession(ctx echo.Context) (*goth.User, bool, error) {
	ouser := &goth.User{}
	b, ok := ctx.Session().Get(OAuthUserSessionKey).([]byte)
	if !ok {
		return nil, ok, ctx.NewErrorWith(ErrGetSession, code.DataStatusIncorrect, ctx.T(`从session中获取oauth2用户信息失败`))
	}
	err := json.Unmarshal(b, ouser)
	return ouser, ok, err
}

func DelSession(ctx echo.Context) {
	ctx.Session().Delete(OAuthUserSessionKey)
}
