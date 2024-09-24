//go:build !bindata
// +build !bindata

package main

import (
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webfront/library/bindata"
)

func initEnv() {
	bindata.Initialize(func() {
		bootconfig.FaviconPath = httpserver.Backend.AssetsDir + `/backend/images/favicon-xs.ico`
	})
}
