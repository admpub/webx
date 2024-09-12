//go:build bindata
// +build bindata

package main

import (
	"github.com/coscms/webcore/library/bindata"
	selfBin "github.com/coscms/webfront/library/bindata"
)

func initEnv() {
	bindata.Asset = Asset
	bindata.AssetDir = AssetDir
	bindata.AssetInfo = AssetInfo
	selfBin.Initialize()
}
