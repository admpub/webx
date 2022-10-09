package oauth2

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func NewConfig(conf *oauth2.Config) *Config {
	return &Config{
		Config: conf,
	}
}

type Config struct {
	*oauth2.Config
}

// Exchange converts an authorization code into a token.
//
// It is used after a resource provider redirects the user back
// to the Redirect URI (the URL obtained from AuthCodeURL).
//
// The HTTP client to use is derived from the context.
// If a client is not provided via the context, http.DefaultClient is used.
//
// The code will be in the *http.Request.FormValue("code"). Before
// calling Exchange, be sure to validate FormValue("state").
func (c *Config) Exchange(ctx context.Context, code interface{}, callbacks ...func(*http.Request)) (*Token, error) {
	var params url.Values
	var callback func(*http.Request)
	switch v := code.(type) {
	case url.Values:
		params = v
	default:
		params = url.Values{
			"grant_type":   {"authorization_code"},
			"code":         {fmt.Sprint(code)},
			"redirect_uri": CondVal(c.RedirectURL),
		}
		if len(c.ClientID) > 0 {
			params.Set("client_id", c.ClientID)
		}
		if len(c.ClientSecret) > 0 {
			params.Set("client_secret", c.ClientSecret)
		}
		/*
			callback = func(req *http.Request) {
				req.SetBasicAuth(url.QueryEscape(c.ClientID), url.QueryEscape(c.ClientSecret))
			}
		*/
	}
	if len(callbacks) > 0 {
		callback = callbacks[0]
	}
	return retrieveToken(ctx, c, params, callback)
}
