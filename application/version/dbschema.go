package version

import (
	"github.com/admpub/nging/v5/application/version"
)

const (
	// 当前应用数据表结构版本号
	dbschema = 0.026
	// 数据表结构最终版本号
	DBSCHEMA = dbschema + version.DBSCHEMA
)
