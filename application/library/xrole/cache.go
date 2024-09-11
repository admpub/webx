package xrole

import (
	"sync"

	"github.com/coscms/webcore/library/perm"
	userNav "github.com/admpub/webx/application/handler/frontend/user/navigate"
	"github.com/webx-top/echo"
)

var (
	navTreeCached *perm.Map
	navTreeOnce   sync.Once
)

func initNavTreeCached() {
	navTreeCached = perm.NewMap(nil)
	navTreeCached.Import(userNav.LeftNavigate)
	navTreeCached.Import(userNav.TopNavigate)
}

func NavTreeCached() *perm.Map {
	navTreeOnce.Do(initNavTreeCached)
	return navTreeCached
}

func init() {
	echo.OnCallback(`nging.httpserver.run.before`, func(_ echo.Event) error {
		NavTreeCached()
		return nil
	})
}
