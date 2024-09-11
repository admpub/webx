package version

import (
	"github.com/coscms/webcore/version"
)

const (
	// 当前应用数据表结构版本号
	dbschema = 0.8
	// 数据表结构最终版本号
	DBSCHEMA = dbschema + version.DBSCHEMA
)
