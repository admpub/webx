//go:build bindata
// +build bindata

package main

import (
	"github.com/admpub/nging/v4/application/library/bindata"
	selfBin "github.com/admpub/webx/application/library/bindata"
)

func initEnv() {
	bindata.Asset = Asset
	bindata.AssetDir = AssetDir
	bindata.AssetInfo = AssetInfo
	selfBin.Initialize()
}
