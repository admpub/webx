package oauth2

import (
	"net/http"
	"net/url"

	"github.com/webx-top/echo"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func NewToken(token *oauth2.Token, raw echo.H) *Token {
	return &Token{
		Token: token,
		Raw:   raw,
	}
}

type Token struct {
	*oauth2.Token

	// raw optionally contains extra metadata from the server
	// when updating a token.
	Raw echo.H
}

// tokenFromInternal maps an *internal.Token struct into
// a *Token struct.
func tokenFromInternal(t *oauth2.Token, raw echo.H) *Token {
	if t == nil {
		return nil
	}
	return NewToken(t, raw)
}

// retrieveToken takes a *Config and uses that to retrieve an *internal.Token.
// This token is then mapped from *internal.Token into an *oauth2.Token which is returned along
// with an error..
func retrieveToken(ctx context.Context, c *Config, v url.Values, cb func(*http.Request)) (*Token, error) {
	tk, raw, err := RetrieveToken(ctx, cb, c.Endpoint.TokenURL, v)
	if err != nil {
		return nil, err
	}
	return tokenFromInternal(tk, raw), nil
}
