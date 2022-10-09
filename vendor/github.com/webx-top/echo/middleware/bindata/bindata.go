/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package bindata

import (
	"net/http"

	"github.com/webx-top/echo"
)

func Static(path string, fs http.FileSystem) echo.MiddlewareFunc {
	length := len(path)
	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			fileName := c.Request().URL().Path()
			if len(fileName) < length || fileName[0:length] != path {
				return next.Handle(c)
			}
			file, err := fs.Open(fileName)
			if err != nil {
				return echo.ErrNotFound
			}
			defer file.Close()
			info, err := file.Stat()
			if err != nil {
				return echo.ErrNotFound
			}
			return c.ServeContent(file, info.Name(), info.ModTime())
		})
	}
}
