// Package alipay implements the OAuth2 protocol for authenticating users through Alipay.
// This package can be used as a reference implementation of an OAuth2 provider for Goth.
package alipay

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/markbates/goth"
	"github.com/smartwalle/crypto4go"

	"github.com/coscms/oauth2s/client/goth/oauth2"
	oauth2x "golang.org/x/oauth2"
)

// These vars define the Authentication, Token, and API URLS for GitHub. If
// using GitHub enterprise you should change these values before calling New.
var (
	SandBoxAuthURL = "https://openauth.alipaydev.com/oauth2/publicAppAuthorize.htm"
	SandBoxAPIURL  = "https://openapi.alipaydev.com/gateway.do"

	AuthURL = "https://openauth.alipay.com/oauth2/publicAppAuthorize.htm"
	APIURL  = "https://openapi.alipay.com/gateway.do"
)

// New creates a new Github provider, and sets up important connection details.
// You should always call `github.New` to get a new Provider. Never try to create
// one manually.
func New(clientKey, secret, callbackURL string, isProduction bool, scopes ...string) *Provider {
	var apiURL, authURL string
	if isProduction {
		apiURL = APIURL
		authURL = AuthURL
	} else {
		apiURL = SandBoxAPIURL
		authURL = SandBoxAuthURL
	}
	return NewCustomisedURL(clientKey, secret, callbackURL, authURL, apiURL, scopes...)
}

// NewCustomisedURL is similar to New(...) but can be used to set custom URLs to connect to
func NewCustomisedURL(clientKey, privateKey, callbackURL, authURL string, apiURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		Secret:       privateKey,
		CallbackURL:  callbackURL,
		HTTPClient:   oauth2.DefaultClient,
		providerName: "alipay",
		profileURL:   apiURL,
	}
	if len(scopes) == 0 {
		scopes = []string{`auth_user`}
	}
	p.config = newConfig(p, authURL, apiURL, scopes)
	return p
}

// Provider is the implementation of `goth.Provider` for accessing Github.
type Provider struct {
	ClientKey     string
	Secret        string
	CallbackURL   string
	HTTPClient    *http.Client
	config        *oauth2.Config
	providerName  string
	profileURL    string
	debug         bool
	appPrivateKey *rsa.PrivateKey // 应用私钥
	once          sync.Once
}

// Name is the name used to retrieve this provider later.
func (p *Provider) Name() string {
	return p.providerName
}

func (p *Provider) parsePrivateKey() error {
	priKey, err := crypto4go.ParsePKCS1PrivateKey(crypto4go.FormatPKCS1PrivateKey(p.Secret))
	if err != nil {
		priKey, err = crypto4go.ParsePKCS8PrivateKey(crypto4go.FormatPKCS8PrivateKey(p.Secret))
		if err != nil {
			return err
		}
	}
	p.appPrivateKey = priKey
	return nil
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

func (p *Provider) createSign(params url.Values) (string, error) {
	var err error
	p.once.Do(func() {
		err = p.parsePrivateKey()
	})
	if err != nil {
		return "", err
	}
	if params.Get(`sign_type`) == `RSA` {
		return SignRSA(params, p.appPrivateKey)
	}
	return SignRSA2(params, p.appPrivateKey)
}

func (p *Provider) urlParams(method string, params url.Values, extra interface{}, scopes ...string) (url.Values, error) {
	if len(method) == 0 {
		method = `alipay.system.oauth.token`
	}
	params.Set("charset", "utf-8")
	params.Set("app_id", p.ClientKey)
	params.Set("method", method)
	params.Set("scope", strings.Join(scopes, `,`))
	params.Set("format", "JSON")
	if extra != nil {
		b, err := json.Marshal(extra)
		if err != nil {
			return nil, err
		}
		params.Set("biz_content", string(b))
	}
	params.Set("sign_type", "RSA2")
	params.Set("timestamp", time.Now().Local().Format(`2006-01-02 15:04:05`))
	params.Set("version", "1.0")
	sign, err := p.createSign(params)
	if err != nil {
		return nil, err
	}
	params.Set("sign", sign)
	return params, nil
}

// Debug is sandbox mode
func (p *Provider) Debug(debug bool) {
	p.debug = debug
}

// BeginAuth asks Github for an authentication end-point.
// documentation https://opensupport.alipay.com/support/helpcenter/166/201602504335?ant_source=zsearch#
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	//url := p.config.AuthCodeURL(state)
	params := url.Values{
		"app_id":       {p.ClientKey},
		"scope":        {"auth_user"},
		"redirect_uri": {p.CallbackURL},
		"state":        {state},
	}
	session := &Session{
		//https://openauth.alipay.com/oauth2/publicAppAuthorize.htm?app_id=商户的APPID&scope=auth_user&redirect_uri=ENCODED_URL&state=init
		AuthURL: p.config.Endpoint.AuthURL + `?` + params.Encode(),
	}
	return session, nil
}

// FetchUser will go to Github and access basic information about the user.
// documentation https://opendocs.alipay.com/apis/api_2/alipay.user.info.share
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*Session)
	user := goth.User{
		AccessToken:  sess.AccessToken,
		RefreshToken: sess.RefreshToken,
		UserID:       sess.OpenID,
		ExpiresAt:    sess.Expiry,
		Provider:     p.Name(),
		RawData:      make(map[string]interface{}),
	}
	if user.AccessToken == "" {
		// data is not yet retrieved since accessToken is still empty
		return user, fmt.Errorf("%s cannot get user information without accessToken", p.providerName)
	}
	param := url.Values{
		"auth_token": {sess.AccessToken},
	}
	var err error
	param, err = p.urlParams(`alipay.user.info.share`, param, nil, `auth_user`)
	if err != nil {
		return user, err
	}
	response, err := p.Client().Get(p.profileURL + `?` + param.Encode())
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return user, fmt.Errorf("Alipay API responded with a %d trying to fetch user information", response.StatusCode)
	}

	bits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	err = json.NewDecoder(bytes.NewReader(bits)).Decode(&user.RawData)
	if err != nil {
		return user, err
	}

	err = userFromReader(bytes.NewReader(bits), &user)
	return user, err
}

const codeSuccess = `10000`

func userFromReader(reader io.Reader, user *goth.User) error {
	/*
		{
		    "alipay_user_info_share_response": {
		        "code": "10000",
		        "msg": "Success",
		        "user_id": "2088102104794936",
		        "avatar": "http://tfsimg.alipay.com/images/partner/T1uIxXXbpXXXXXXXX",
		        "province": "安徽省",
		        "city": "安庆",
		        "nick_name": "支付宝小二",
		        "gender": "F"
		    },
		    "sign": "ERITJKEIJKJHKKKKKKKHJEREEEEEEEEEEE"
		}
	*/
	u := struct {
		Response struct {
			Code      string `json:"code"`
			Msg       string `json:"msg"`
			Name      string `json:"nick_name"`
			AvatarURL string `json:"avatar"`
			Gender    string `json:"gender"`
			Province  string `json:"province"` // 省份名称
			City      string `json:"city"`     // 城市名称
			UserID    string `json:"user_id"`  // 支付宝用户userId 最大长度16位
		} `json:"alipay_user_info_share_response"`
		Sign string `json:"sign"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}
	if u.Response.Code != codeSuccess {
		return errors.New(u.Response.Msg)
	}

	user.Name = u.Response.Name
	user.NickName = u.Response.Name
	user.AvatarURL = u.Response.AvatarURL
	user.RawData[`gender`] = u.Response.Gender
	if len(u.Response.Province) > 0 && len(u.Response.City) > 0 {
		user.Location = u.Response.Province + `,` + u.Response.City
	}
	user.IDToken = u.Response.UserID
	if len(u.Response.UserID) > 0 {
		user.UserID = u.Response.UserID
	}

	return err
}

func newConfig(provider *Provider, authURL, tokenURL string, scopes []string) *oauth2.Config {
	c := &oauth2x.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2x.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	for _, scope := range scopes {
		c.Scopes = append(c.Scopes, scope)
	}

	return oauth2.NewConfig(c)
}

//RefreshToken refresh token is not provided by QQ
func (p *Provider) RefreshToken(refreshToken string) (*oauth2x.Token, error) {
	return nil, errors.New("Refresh token is not provided by alipay")
}

//RefreshTokenAvailable refresh token is not provided by QQ
func (p *Provider) RefreshTokenAvailable() bool {
	return false
}
