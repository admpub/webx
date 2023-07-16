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

// 也可以以服务的方式启动nging
// 服务支持的操作有：
// nging service install  	-- 安装服务
// nging service uninstall  -- 卸载服务
// nging service start 		-- 启动服务
// nging service stop 		-- 停止服务
// nging service restart 	-- 重启服务
package main

//go:generate go install github.com/admpub/bindata/v3/go-bindata@latest
//go:generate go-bindata -fs -o bindata_assetfs.go -prefix "../nging/|vendor/github.com/nging-plugins/dbmanager/" -ignore "\\.(DS_Store|less|scss)$" -minify "\\.(js|css)$" -tags bindata ../nging/public/assets/... ../nging/template/... public/assets/... template/... config/i18n/... vendor/github.com/nging-plugins/dbmanager/template/... vendor/github.com/nging-plugins/dbmanager/public/assets/...

import (
	"os"
	"time"

	_ "github.com/admpub/bindata/v3"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/cmd"
	"github.com/admpub/nging/v5/application/cmd/bootconfig"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/license"
	"github.com/admpub/nging/v5/application/library/module"

	// register

	_ "github.com/admpub/nging/v5/application/initialize/manager" // must first

	_ "github.com/admpub/webx/application"
	"github.com/admpub/webx/application/initialize/initnavigate"

	//_ "github.com/nging-plugins/open/application/library/oauth2server/initexample"

	// swagger docs
	//_ "github.com/nging-plugins/open/application/handler/frontend/swagger"

	"github.com/admpub/webx/application/handler/frontend/article"
	"github.com/admpub/webx/application/handler/frontend/index"
	"github.com/admpub/webx/application/version"

	// module
	"github.com/admpub/nging/v5/application/handler/cloud"
	"github.com/admpub/nging/v5/application/handler/task"
	"github.com/nging-plugins/dbmanager"
)

var (
	BUILD_TIME string
	BUILD_OS   string
	BUILD_ARCH string
	CLOUD_GOX  string
	COMMIT     string
	LABEL      = `dev` //beta/alpha/stable
	VERSION    = `1.0.0`
	PACKAGE    = `free`

	schemaVer = version.DBSCHEMA //数据表结构版本
	name      = `Webx`
)

func main() {
	defer log.Close()
	index.DefaultIndexHandler = article.Index
	index.DefaultSearchHandler = article.List

	echo.Set(`BUILD_TIME`, BUILD_TIME)
	echo.Set(`BUILD_OS`, BUILD_OS)
	echo.Set(`BUILD_ARCH`, BUILD_ARCH)
	echo.Set(`COMMIT`, COMMIT)
	echo.Set(`LABEL`, LABEL)
	echo.Set(`VERSION`, VERSION)
	echo.Set(`PACKAGE`, PACKAGE)
	echo.Set(`SCHEMA_VER`, schemaVer)
	config.Version.Name = name
	bootconfig.SoftwareName = name
	bootconfig.Short = name
	bootconfig.Long = ``
	bootconfig.Welcome = "Thank you for choosing " + name + " %s, I hope you enjoy using it.\nToday is %s."
	if com.FileExists(`config/install.sql`) {
		os.Rename(`config/install.sql`, `config/install.sql.`+time.Now().Format(`20060102150405.000`))
	}
	initModule()
	initEnv()
	exec()
}

func exec() {
	cmd.Execute()
}

func initModule() {
	module.Register(
		&task.Module,
		&cloud.Module,
		&dbmanager.Module,
	)
	initnavigate.Initialize()
}

func init() {
	license.SetProductName(`webx`)
}
