// Package qq implements the OAuth2 protocol for authenticating users through QQ.
// This package can be used as a reference implementation of an OAuth2 provider for Goth.
package qq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/coscms/oauth2s/client/goth/oauth2"
	"github.com/markbates/goth"
	oauth2x "golang.org/x/oauth2"
)

// These vars define the Authentication, Token, and API URLS for GitHub. If
// using GitHub enterprise you should change these values before calling New.
var (
	AuthURL    = "https://graph.qq.com/oauth2.0/authorize"
	TokenURL   = "https://graph.qq.com/oauth2.0/token"
	ProfileURL = "https://graph.qq.com/user/get_user_info"
	MeURL      = "https://graph.qq.com/oauth2.0/me"
)

// New creates a new Github provider, and sets up important connection details.
// You should always call `github.New` to get a new Provider. Never try to create
// one manually.
func New(clientKey, secret, callbackURL string, scopes ...string) *Provider {
	return NewCustomisedURL(clientKey, secret, callbackURL, AuthURL, TokenURL, ProfileURL, MeURL, scopes...)
}

// NewCustomisedURL is similar to New(...) but can be used to set custom URLs to connect to
func NewCustomisedURL(clientKey, secret, callbackURL, authURL, tokenURL, profileURL, meURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		Secret:       secret,
		CallbackURL:  callbackURL,
		HTTPClient:   oauth2.DefaultClient,
		providerName: "qq",
		profileURL:   profileURL,
		meURL:        meURL,
	}
	if len(scopes) == 0 {
		scopes = []string{`get_user_info`, `add_share`}
	}
	p.config = newConfig(p, authURL, tokenURL, scopes)
	return p
}

// Provider is the implementation of `goth.Provider` for accessing Github.
type Provider struct {
	ClientKey    string
	Secret       string
	CallbackURL  string
	HTTPClient   *http.Client
	config       *oauth2x.Config
	providerName string
	profileURL   string
	meURL        string
}

// Name is the name used to retrieve this provider later.
func (p *Provider) Name() string {
	return p.providerName
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

func (p *Provider) urlParams(sess *Session) string {
	return `?oauth_consumer_key=` + url.QueryEscape(p.ClientKey) + `&access_token=` + url.QueryEscape(sess.AccessToken) + `&openid=` + url.QueryEscape(sess.OpenID) + `&format=json`
}

// Debug is a no-op for the github package.
func (p *Provider) Debug(debug bool) {}

// BeginAuth asks Github for an authentication end-point.
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	session := &Session{
		AuthURL: p.config.AuthCodeURL(state),
	}
	return session, nil
}

// FetchUser will go to Github and access basic information about the user.
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*Session)
	user := goth.User{
		AccessToken:  sess.AccessToken,
		RefreshToken: sess.RefreshToken,
		ExpiresAt:    sess.Expiry,
		Provider:     p.Name(),
	}

	if user.AccessToken == "" {
		// data is not yet retrieved since accessToken is still empty
		return user, fmt.Errorf("%s cannot get user information without accessToken", p.providerName)
	}
	if err := getOpenID(p, sess); err != nil {
		return user, err
	}

	user.UserID = sess.OpenID

	response, err := p.Client().Get(p.profileURL + p.urlParams(sess))
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return user, fmt.Errorf("QQ API responded with a %d trying to fetch user information", response.StatusCode)
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
	if err != nil {
		return user, err
	}

	user.RawData[`unionid`] = sess.UnionID
	return user, err
}

func userFromReader(reader io.Reader, user *goth.User) error {
	u := struct {
		Name      string `json:"nickname"`
		AvatarURL string `json:"figureurl_2"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Name = u.Name
	user.NickName = u.Name
	//user.Email = u.Email
	user.AvatarURL = u.AvatarURL
	//user.UserID = strconv.Itoa(u.ID)
	//user.Location = u.Location
	return err
}

func getOpenID(p *Provider, sess *Session) error {
	if len(sess.OpenID) > 0 {
		return nil
	}
	if len(sess.AccessToken) == 0 {
		return fmt.Errorf("%s cannot get openid without accessToken", p.providerName)
	}
	response, err := p.Client().Get(p.meURL + `?access_token=` + url.QueryEscape(sess.AccessToken))
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("QQ API responded with a %d trying to fetch user openid", response.StatusCode)
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	b = bytes.TrimSpace(b)
	b = bytes.TrimPrefix(b, []byte(`callback(`))
	b = bytes.TrimSuffix(b, []byte(`};`))
	data := map[string]string{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	var ok bool
	sess.OpenID, ok = data["openid"]
	if !ok || len(sess.OpenID) == 0 {
		if description, ok := data["error_description"]; ok {
			err = errors.New(description)
		} else {
			err = errors.New(`Cannot get openid`)
		}
	}
	return err
}

func newConfig(provider *Provider, authURL, tokenURL string, scopes []string) *oauth2x.Config {
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

	return c
}

//RefreshToken refresh token is not provided by QQ
func (p *Provider) RefreshToken(refreshToken string) (*oauth2x.Token, error) {
	return nil, errors.New("Refresh token is not provided by QQ")
}

//RefreshTokenAvailable refresh token is not provided by QQ
func (p *Provider) RefreshTokenAvailable() bool {
	return false
}
