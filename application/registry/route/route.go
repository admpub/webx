/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package route

import (
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/route"
)

var (
	routeRegister = route.NewRegister(echo.New())
)

func init() {
	route.Default.Frontend = routeRegister
}

func IRegister() route.IRegister {
	return routeRegister
}

func MetaHandler(handler interface{}, m echo.H, requests ...interface{}) echo.Handler {
	return routeRegister.MetaHandler(m, handler, requests...)
}

func MakeHandler(handler interface{}, requests ...interface{}) echo.Handler {
	return routeRegister.MakeHandler(handler, requests...)
}

func MetaHandlerWithRequest(handler interface{}, m echo.H, requests interface{}, methods ...string) echo.Handler {
	return routeRegister.MetaHandlerWithRequest(m, handler, requests, methods...)
}

func HandlerWithRequest(handler interface{}, requests interface{}, methods ...string) echo.Handler {
	return routeRegister.HandlerWithRequest(handler, requests, methods...)
}

func Pre(middlewares ...interface{}) {
	routeRegister.Pre(middlewares...)
}

func PreToGroup(groupName string, middlewares ...interface{}) {
	routeRegister.PreToGroup(groupName, middlewares...)
}

func Use(middlewares ...interface{}) {
	routeRegister.Use(middlewares...)
}

func UseToGroup(groupName string, middlewares ...interface{}) {
	routeRegister.UseToGroup(groupName, middlewares...)
}

func Host(hostName string, middlewares ...interface{}) route.Hoster {
	return routeRegister.Host(hostName, middlewares...)
}

func Register(fn func(echo.RouteRegister)) {
	routeRegister.Register(fn)
}

func RegisterToGroup(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
	routeRegister.RegisterToGroup(groupName, fn, middlewares...)
}

func Apply() {
	echo.PanicIf(echo.Fire(`webx.route.apply.before`))
	routeRegister.Apply()
	echo.PanicIf(echo.Fire(`webx.route.apply.after`))
}
