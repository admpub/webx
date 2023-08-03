package middleware

import (
	"regexp"

	"github.com/webx-top/echo"
)

func HostChecker() echo.MiddlewareFuncd {
	key := `frontend.hostRuleRegexp`
	return func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			re, ok := echo.Get(key).(*regexp.Regexp)
			if !ok {
				return h.Handle(c)
			}
			// c.Host() 不含端口号
			if re.MatchString(c.Host()) {
				return h.Handle(c)
			}
			c.Response().NotFound()
			return nil
		}
	}
}
