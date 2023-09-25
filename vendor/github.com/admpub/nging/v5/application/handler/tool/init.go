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

package tool

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/tool`, func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/ip`, IP2Region)
		g.Route(`GET,POST`, `/base64`, Base64)
		g.Route(`GET,POST`, `/url`, URL)
		g.Route(`GET,POST`, `/timestamp`, Timestamp)
		g.Route(`GET,POST`, `/regexp_test`, RegexpTest)
		g.Route(`GET,POST`, `/replaceurl`, ReplaceURL)
		g.Route(`GET,POST`, `/gen_password`, GenPassword)
	})
}
