package mpwechat

import (
	"errors"
	"io"
	"net/url"
	"time"

	"github.com/admpub/goth"
	"github.com/admpub/marmot/miner"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"
)

type WechatAuthResult struct {
	AccessToken  string `json:"access_token"`  //: "39_pW95pgvNPYWl6cwnTa_****************************************************-mw0cpC85m1AJ16HCsh0O-Z-ZVDigLsoJGFQ",
	Code         string `json:"code"`          //: "0232DTkl2F2j564O5qnl2TAhdO32DTkq",
	ExpiresIn    int64  `json:"expires_in"`    //: 7200,
	OpenID       string `json:"openid"`        //: "oRrdQt8*************QvsqGu_U",
	RefreshToken string `json:"refresh_token"` //: "39_H4WOg6FqOtktcLSiwclqYaNEoDnnUA-jPPfCx769isAstsp77Z64U7FFy5Com8dqN58iqWnilZgF0lyDU16lEn8_Ti89UaU_*******-***",
	Scope        string `json:"scope"`         //: "snsapi_userinfo",
	UnionID      string `json:"unionid"`       //: "oU5Yyt48Fmo*********_HR5D0vc"
}

type WechatUserInfo struct {
	OpenID    string `json:"openId"`    //: "oRrdQt8*************QvsqGu_U",
	NickName  string `json:"nickName"`  //: "辉煌",
	Gender    string `json:"gender"`    //: 1,
	City      string `json:"city"`      //: "Guangzhou",
	Province  string `json:"province"`  //: "Guangdong",
	Country   string `json:"country"`   //: "China",
	AvatarURL string `json:"avatarUrl"` //: "https://thirdwx.qlogo.cn/mmopen/vi_32/PWLwicVtiaRhHQNTDpGIh9iauwCPQ4uI7H1Q2nGoicSd4VTxqoMjmBhjdPO1YzFlFic6RkUFmGM54k6I2gaNGIGwGMg/132",
	UnionID   string `json:"unionId"`   //: "oU5Yyt48Fmo*********_HR5D0vc"
}

func (u *WechatUserInfo) Location() string {
	if len(u.Country) == 0 && len(u.Province) == 0 && len(u.City) == 0 {
		return ``
	}
	return u.Country + `, ` + u.Province + `, ` + u.City
}

type WechatPostData struct {
	AuthResult *WechatAuthResult `json:"authResult"`
	UserInfo   *WechatUserInfo   `json:"userInfo"`
}

func NewWechatPostData() *WechatPostData {
	return &WechatPostData{
		AuthResult: &WechatAuthResult{},
		UserInfo:   &WechatUserInfo{},
	}
}

func (post *WechatPostData) Check(ctx echo.Context) error {
	if post.AuthResult.OpenID != post.UserInfo.OpenID {
		return ctx.NewError(stdCode.InvalidParameter, ctx.T(`openId不一致`)).SetZone(`openid`)
	}
	if post.AuthResult.UnionID != post.UserInfo.UnionID {
		return ctx.NewError(stdCode.InvalidParameter, ctx.T(`unionId不一致`)).SetZone(`unionid`)
	}
	return nil
}

var APIURL = `https://api.weixin.qq.com/sns/jscode2session`

func (post *WechatPostData) Post(ctx echo.Context, appID string, appSecret string) (*Response, error) {
	apiURL := APIURL
	apiValues := url.Values{}
	apiValues.Set(`appid`, appID)
	apiValues.Set(`secret`, appSecret)
	apiValues.Set(`js_code`, post.AuthResult.Code)
	apiValues.Set(`grant_type`, `authorization_code`)
	client := miner.NewClient(time.Second * 10)
	resp, err := client.Get(apiURL + `?` + apiValues.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := &Response{}
	err = com.JSONDecode(b, result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode > 0 {
		return nil, errors.New(result.ErrMsg)
	}
	if post.AuthResult.UnionID != result.UnionID {
		return nil, ctx.NewError(stdCode.InvalidParameter, ctx.T(`unionId不正确`)).SetZone(`unionid`)
	}
	if post.AuthResult.OpenID != result.OpenID {
		return nil, ctx.NewError(stdCode.InvalidParameter, ctx.T(`openId不正确`)).SetZone(`openid`)
	}
	return result, nil
}

type Response struct {
	// 错误时结果
	ErrCode int64  `json:"errcode,omitempty" xml:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty" xml:"errmsg,omitempty"`

	// 正常时结果
	UnionID    string `json:"unionid,omitempty" xml:"unionid,omitempty"`
	ExpiresIn  int64  `json:"expires_in,omitempty" xml:"expires_in,omitempty"`   //凭证有效时间，单位：秒
	OpenID     string `json:"openid,omitempty" xml:"openid,omitempty"`           //用户唯一标识
	SessionKey string `json:"session_key,omitempty" xml:"session_key,omitempty"` //会话密匙（考虑到应用安全，不应该在网络上传输session_key）
}

func (result *Response) AsUser(post *WechatPostData) *goth.User {
	return &goth.User{
		Provider:          `wechat`,
		UserID:            result.OpenID,
		ExpiresAt:         time.Now().Add(time.Second * time.Duration(result.ExpiresIn)),
		Email:             ``,
		Name:              ``,
		FirstName:         ``,
		LastName:          ``,
		NickName:          post.UserInfo.NickName,
		AvatarURL:         post.UserInfo.AvatarURL,
		Location:          post.UserInfo.Location(),
		AccessToken:       post.AuthResult.AccessToken,
		AccessTokenSecret: ``,
		RefreshToken:      post.AuthResult.RefreshToken,
		IDToken:           ``,
		RawData: map[string]interface{}{
			`unionid`: result.UnionID,
			`gender`:  post.UserInfo.Gender,
		},
	}
}
