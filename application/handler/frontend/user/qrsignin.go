package user

import (
	"math"
	"sync/atomic"
	"time"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webfront/library/cache"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

type QRSignIn struct {
	SessionID     string `json:"sID"`
	SessionMaxAge int    `json:"sAge"`
	Expires       int64  `json:"dExp"` // 数据过期时间戳(秒)
	IPAddress     string `json:"ip"`
	UserAgent     string `json:"ua"`
	Platform      string `json:"pf"`
	Scense        string `json:"ss"`
	DeviceNo      string `json:"dn"`
}

var autoIncr atomic.Uint64

func GetMaxNumber() uint64 {
	v := autoIncr.Add(1)
	if v >= math.MaxUint64-1 {
		autoIncr.Store(0)
	}
	return v
}

func (q QRSignIn) Encode() (string, error) {
	plaintext, err := com.JSONEncodeToString(q)
	if err != nil {
		return ``, err
	}
	qrcode := config.FromFile().Encode256(plaintext)
	qrcode = com.URLSafeBase64(qrcode, true)
	return qrcode, nil
}

func (q *QRSignIn) Decode(ctx echo.Context, encrypted string) error {
	encrypted = com.URLSafeBase64(encrypted, false)
	plaintext := config.FromFile().Decode256(encrypted)
	if len(plaintext) == 0 {
		return ctx.NewError(code.InvalidParameter, `解密失败`).SetZone(`data`)
	}
	return com.JSONDecodeString(plaintext, q)
}

func NewQRSignIn(ctx echo.Context, cookieMaxAge int, expireTime time.Time) QRSignIn {
	qsi := QRSignIn{
		SessionID:     ctx.Session().MustID(),
		SessionMaxAge: cookieMaxAge,
		Expires:       expireTime.Unix(),
		IPAddress:     ctx.RealIP(),
		UserAgent:     ctx.Request().UserAgent(),
		Platform:      ctx.Header(`X-Platform`),
		Scense:        ctx.Header(`X-Scense`),
		DeviceNo:      ctx.Header(`X-Device-Id`),
	}
	if len(qsi.Scense) > 0 {
		qsi.Scense = `qrcode_` + qsi.Scense
	} else {
		qsi.Scense = `qrcode_` + modelCustomer.DefaultDeviceScense
	}
	return qsi
}

func GenerateUniqueKey(ip, ua string) string {
	return com.Md5(com.String(time.Now().UnixMicro())+ua+ip) + com.RandomAlphanumeric(2)
}

type QRSignInCase interface {
	Encode(echo.Context, QRSignIn) (string, error)
	Decode(echo.Context, string) (QRSignIn, error)
}

type cacheQRSignIn struct {
}

func (c cacheQRSignIn) Encode(ctx echo.Context, signInData QRSignIn) (string, error) {
	key := GenerateUniqueKey(ctx.RealIP(), ctx.Request().UserAgent())
	err := cache.Put(ctx, `qrsignin_`+key, signInData, signInData.Expires-time.Now().Unix())
	return key, err
}

func (c cacheQRSignIn) Decode(ctx echo.Context, key string) (QRSignIn, error) {
	signInData := QRSignIn{}
	if !com.StrIsAlphaNumeric(key) {
		return signInData, ctx.NewError(code.InvalidParameter, `无效字符`).SetZone(`data`)
	}
	err := cache.Get(ctx, `qrsignin_`+key, &signInData)
	return signInData, err
}

type defaultQRSignIn struct {
}

func (c defaultQRSignIn) Encode(_ echo.Context, signInData QRSignIn) (string, error) {
	return signInData.Encode()
}

func (c defaultQRSignIn) Decode(ctx echo.Context, encrypted string) (QRSignIn, error) {
	signInData := QRSignIn{}
	err := signInData.Decode(ctx, encrypted)
	return signInData, err
}

var qrSignInCases = map[string]QRSignInCase{
	`cache`:   cacheQRSignIn{},   // 缺点：占用存储空间；优点：字符串短，生成的二维码更容易识别
	`default`: defaultQRSignIn{}, // 优点：不占用存储空间；缺点：加密字符串太长，生成的二维码元素图块小而多，不易识别
}

func GetQRSignInCase(caseName string) QRSignInCase {
	cs, ok := qrSignInCases[caseName]
	if ok {
		return cs
	}
	return qrSignInCases[`default`]
}

func RegisterQRSignInCase(caseName string, qrsic QRSignInCase) {
	qrSignInCases[caseName] = qrsic
}
