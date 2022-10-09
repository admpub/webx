//go:build !bindata
// +build !bindata

package main

import (
	"github.com/admpub/nging/v4/application/cmd/bootconfig"
	"github.com/admpub/nging/v4/application/initialize/backend"
	"github.com/admpub/webx/application/library/bindata"
)

func initEnv() {
	bindata.Initialize(func() {
		bootconfig.FaviconPath = backend.AssetsDir + `/backend/images/favicon-xs.ico`
	})
}
