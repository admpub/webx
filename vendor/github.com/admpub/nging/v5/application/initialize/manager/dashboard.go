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

package manager

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/nging/v5/application/registry/dashboard"
)

func init() {
	dashboard.CardRegister(
		(&dashboard.Card{
			IconName:  `fa-tasks`,
			IconColor: `danger`,
			Short:     `TASKS`,
			Name:      `计划任务数量`,
			Summary:   ``,
		}).SetContentGenerator(func(ctx echo.Context) interface{} {
			//计划任务统计
			taskMdl := model.NewTask(ctx)
			taskCount, _ := taskMdl.Count(nil)
			return taskCount
		}),
	)

}
