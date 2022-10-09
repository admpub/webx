package oauth2

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
)

func CondVal(v string) []string {
	if len(v) == 0 {
		return nil
	}
	return []string{v}
}

// HTTPClient is the context key to use with golang.org/x/net/context's
// WithValue function to associate an *http.Client value with a context.
var HTTPClient ContextKey

// ContextKey is just an empty struct. It exists so HTTPClient can be
// an immutable public variable with a unique type. It's immutable
// because nobody else can create a ContextKey, being unexported.
type ContextKey struct{}

// ContextClientFunc is a func which tries to return an *http.Client
// given a Context value. If it returns an error, the search stops
// with that error.  If it returns (nil, nil), the search continues
// down the list of registered funcs.
type ContextClientFunc func(context.Context) (*http.Client, error)

var contextClientFuncs []ContextClientFunc

func RegisterContextClientFunc(fn ContextClientFunc) {
	contextClientFuncs = append(contextClientFuncs, fn)
}

func ContextClient(ctx context.Context) (*http.Client, error) {
	if ctx != nil {
		if hc, ok := ctx.Value(HTTPClient).(*http.Client); ok {
			return hc, nil
		}
	}
	for _, fn := range contextClientFuncs {
		c, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		if c != nil {
			return c, nil
		}
	}
	return http.DefaultClient, nil
}

func RetrieveToken(ctx context.Context, callback func(req *http.Request), tokenURL string, v url.Values) (*oauth2.Token, echo.H, error) {
	var token *oauth2.Token
	raw := echo.H{}
	hc, err := ContextClient(ctx)
	if err != nil {
		return token, raw, err
	}
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return token, raw, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if callback != nil {
		callback(req)
	}
	r, err := ctxhttp.Do(ctx, hc, req)
	if err != nil {
		return token, raw, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return token, raw, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}
	if code := r.StatusCode; code < 200 || code > 299 {
		return token, raw, fmt.Errorf("oauth2: cannot fetch token: %v\nResponse: %s", r.Status, body)
	}

	var lifetime int64
	content, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch content {
	case "application/x-www-form-urlencoded", "text/plain":
		vals, err := url.ParseQuery(string(body))
		if err != nil {
			return token, raw, err
		}
		token = &oauth2.Token{
			AccessToken:  vals.Get("access_token"),
			TokenType:    vals.Get("token_type"),
			RefreshToken: vals.Get("refresh_token"),
		}

		for k, v := range vals {
			if len(v) == 0 {
				continue
			}
			if len(v) == 1 {
				raw[k] = v[0]
			} else {
				raw[k] = v
			}
		}

		expires := vals.Get("expires_in")
		if len(expires) == 0 {
			expires = vals.Get("expires")
		}
		if len(expires) > 0 {
			lifetime = param.AsInt64(expires)
		}
	default:
		token = &oauth2.Token{}
		err = json.Unmarshal(body, &raw)
		if err != nil {
			return token, raw, err
		}
		var m echo.H
		if !raw.Has(`access_token`) {
			for _, val := range raw {
				v, y := val.(map[string]interface{})
				if !y {
					continue
				}
				_, y = v[`access_token`]
				if !y {
					continue
				}
				m = echo.H(v)
				break
			}
		} else {
			m = raw
		}
		if m != nil {
			token.AccessToken = m.String(`access_token`)
			token.RefreshToken = m.String(`refresh_token`)
			token.TokenType = m.String(`token_type`)
			if m.Has(`expires_in`) {
				lifetime = m.Int64(`expires_in`)
			} else if m.Has(`expires`) {
				lifetime = m.Int64(`expires`)
			}
		}
	}
	if lifetime > 0 {
		token.Expiry = time.Now().Local().Add(time.Duration(lifetime) * time.Second)
	}
	// Don't overwrite `RefreshToken` with an empty value
	// if this was a token refreshing request.
	if len(token.RefreshToken) == 0 {
		token.RefreshToken = v.Get("refresh_token")
	}
	return token, raw, nil
}
